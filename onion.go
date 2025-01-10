package onion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc defines the function signature for route handlers.
type HandlerFunc func(*Context)

// Context wraps http.ResponseWriter and *http.Request, plus path parameters.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	params   map[string]string
}

// String is a helper for sending plain text.
func (c *Context) String(statusCode int, msg string) {
	c.Response.WriteHeader(statusCode)
	c.Response.Write([]byte(msg))
}

// JSON is a helper for sending JSON data.
func (c *Context) JSON(statusCode int, data interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(statusCode)
	json.NewEncoder(c.Response).Encode(data)
}

// Param fetches a path param like ":bookId".
func (c *Context) Param(key string) string {
	return c.params[key]
}

// ----------------------------------------------------
// App (the main Onion application struct)
// ----------------------------------------------------

type App struct {
	mux         *http.ServeMux
	middlewares []HandlerFunc
	notFound    HandlerFunc

	// We'll store routes here in a map, keyed by (method, pattern)
	routes map[routeKey]HandlerFunc
}

type routeKey struct {
	method  string
	pattern string
}

// Route defines a single HTTP route.
type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
}

// New creates a new Onion app
func New() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []HandlerFunc{},
		notFound: func(c *Context) {
			http.NotFound(c.Response, c.Request)
		},
		routes: make(map[routeKey]HandlerFunc),
	}
}

// Use registers a middleware that will run before route handlers.
func (a *App) Use(mw HandlerFunc) {
	a.middlewares = append(a.middlewares, mw)
}

// NotFoundHandler sets a custom 404.
func (a *App) NotFoundHandler(fn HandlerFunc) {
	a.notFound = fn
}

// UseRoutes loads multiple route slices (like BookRoutes, UserRoutes).
func (a *App) UseRoutes(routeGroups ...[]Route) {
	for _, group := range routeGroups {
		for _, r := range group {
			a.handle(r.Method, r.Pattern, r.Handler)
		}
	}
}

// handle just stores the route in our map. We do the actual matching in dispatch().
func (a *App) handle(method, pattern string, handler HandlerFunc) {
	a.routes[routeKey{method, pattern}] = handler
}

// Run starts the server. Here we register one wildcard handleFunc to dispatch.
func (a *App) Run(addr string) error {
	fmt.Println("Onion server running on", addr)

	// Register exactly one fallback route: "/"
	a.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.dispatch(w, r)
	})

	return http.ListenAndServe(addr, a.mux)
}

// dispatch finds a matching route by (method, path), extracts params, executes middlewares, etc.
func (a *App) dispatch(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	reqMethod := r.Method

	// We'll do a param-capable match. For example, if the user route is "/books/:bookId"
	// and the incoming path is "/books/123", we want to pick that route and fill param "bookId" = "123".
	//
	// Steps:
	//   1) Scan all known routes for any that match the method
	//   2) For each route with same method, check if the path matches (with param placeholders)
	//   3) If found, parse out params and call its handler
	//   4) Otherwise fallback to 404

	for key, handler := range a.routes {
		if key.method == reqMethod {
			params, ok := matchWithParams(key.pattern, reqPath)
			if ok {
				c := &Context{
					Response: w,
					Request:  r,
					params:   params,
				}

				// Middlewares
				for _, mw := range a.middlewares {
					mw(c)
				}

				// Handler
				handler(c)
				return
			}
		}
	}

	// If we reach here, no route matched => 404
	a.notFound(&Context{Response: w, Request: r})
}

// matchWithParams checks if the "pattern" (like "/books/:bookId") matches "path" ("/books/123").
// If it matches, returns (map[string]string, true). If not, returns (nil, false).
func matchWithParams(pattern, path string) (map[string]string, bool) {
	pParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	// They must have the same number of segments
	if len(pParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i := 0; i < len(pParts); i++ {
		pp := pParts[i]
		pa := pathParts[i]

		if strings.HasPrefix(pp, ":") {
			// param placeholder
			key := strings.TrimPrefix(pp, ":")
			params[key] = pa
		} else if pp != pa {
			// mismatch
			return nil, false
		}
	}

	return params, true
}

// ----------------------------------------------------
// RouteGroup (Fluent group builder)
// ----------------------------------------------------

type RouteGroup struct {
	prefix string
	routes []Route
}

// NewGroup("books") => prefix = "books"
func NewGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix: prefix,
		routes: []Route{},
	}
}

// GET etc. Just appends a Route with the correct method, path, handler
func (rg *RouteGroup) GET(pattern string, handler HandlerFunc) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:  http.MethodGet,
		Pattern: "/" + rg.prefix + pattern,
		Handler: handler,
	})
	return rg
}

func (rg *RouteGroup) POST(pattern string, handler HandlerFunc) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:  http.MethodPost,
		Pattern: "/" + rg.prefix + pattern,
		Handler: handler,
	})
	return rg
}

func (rg *RouteGroup) PUT(pattern string, handler HandlerFunc) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:  http.MethodPut,
		Pattern: "/" + rg.prefix + pattern,
		Handler: handler,
	})
	return rg
}

func (rg *RouteGroup) DELETE(pattern string, handler HandlerFunc) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:  http.MethodDelete,
		Pattern: "/" + rg.prefix + pattern,
		Handler: handler,
	})
	return rg
}

// Routes returns the final []Route
func (rg *RouteGroup) Routes() []Route {
	return rg.routes
}
