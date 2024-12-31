package onion

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc defines the function signature for route handlers
type HandlerFunc func(*Context)

// Context wraps http.ResponseWriter and *http.Request, plus path parameters and utilities
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	params   map[string]string
}

// String sends a plain text response
func (c *Context) String(statusCode int, msg string) {
	c.Response.WriteHeader(statusCode)
	c.Response.Write([]byte(msg))
}

// JSON sends a JSON response
func (c *Context) JSON(statusCode int, data interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(statusCode)
	fmt.Fprintf(c.Response, "%v", data) // Simple JSON encoding (replace with json.Marshal if needed)
}

// Param retrieves a path parameter by name (e.g., "/users/:id" -> "id")
func (c *Context) Param(key string) string {
	return c.params[key]
}

// ----------------------------------------------------
// App: The main Onion application struct
// ----------------------------------------------------
type App struct {
	mux         *http.ServeMux
	middlewares []HandlerFunc
	notFound    HandlerFunc
}

// Route defines a single HTTP route
type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
}

// New creates a new Onion application
func New() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []HandlerFunc{},
		notFound: func(c *Context) {
			http.NotFound(c.Response, c.Request)
		},
	}
}

// Use adds a global middleware
func (a *App) Use(mw HandlerFunc) {
	a.middlewares = append(a.middlewares, mw)
}

// NotFoundHandler sets a custom 404 handler
func (a *App) NotFoundHandler(fn HandlerFunc) {
	a.notFound = fn
}

// MapRoutes maps multiple route slices to the application
func (a *App) MapRoutes(routeGroups ...[]Route) {
	for _, groupRoutes := range routeGroups {
		for _, r := range groupRoutes {
			a.handle(r.Method, r.Pattern, r.Handler)
		}
	}
}

// handle registers a single route with the application
func (a *App) handle(method, pattern string, handler HandlerFunc) {
	a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			a.notFound(&Context{Response: w, Request: r})
			return
		}

		// Extract path parameters
		params := extractParams(pattern, r.URL.Path)

		// Prepare the context
		c := &Context{
			Response: w,
			Request:  r,
			params:   params,
		}

		// Execute middlewares
		for _, mw := range a.middlewares {
			mw(c)
		}

		// Call the route handler
		handler(c)
	})
}

// Run starts the Onion server
func (a *App) Run(addr string) error {
	fmt.Printf("Onion server running on %s\n", addr)
	return http.ListenAndServe(addr, a.mux)
}

// ----------------------------------------------------
// Helper Functions
// ----------------------------------------------------

// extractParams extracts path parameters based on the URL pattern
func extractParams(pattern, path string) map[string]string {
	params := make(map[string]string)
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	for i := 0; i < len(patternParts) && i < len(pathParts); i++ {
		if strings.HasPrefix(patternParts[i], ":") {
			key := strings.TrimPrefix(patternParts[i], ":")
			params[key] = pathParts[i]
		}
	}

	return params
}
