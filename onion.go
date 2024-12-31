package onion

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc defines the function signature for route handlers.
// Each handler takes a pointer to Context, which provides convenience methods
// for working with HTTP requests and responses.
type HandlerFunc func(*Context)

// Context wraps http.ResponseWriter and *http.Request, plus path parameters and utilities.
// This struct provides helper methods for sending responses and accessing request data.
type Context struct {
	Response http.ResponseWriter // The response writer for sending HTTP responses
	Request  *http.Request       // The incoming HTTP request
	params   map[string]string   // Path parameters extracted from the URL
}

// String sends a plain text response with a given status code.
func (c *Context) String(statusCode int, msg string) {
	c.Response.WriteHeader(statusCode)
	c.Response.Write([]byte(msg))
}

// JSON sends a JSON response with a given status code.
// For simplicity, it uses fmt.Fprintf to serialize the data.
func (c *Context) JSON(statusCode int, data interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(statusCode)
	fmt.Fprintf(c.Response, "%v", data) // Replace with json.Marshal for proper JSON serialization
}

// Param retrieves a path parameter by name.
// Example: Given the route "/users/:id" and path "/users/123",
// calling c.Param("id") will return "123".
func (c *Context) Param(key string) string {
	return c.params[key]
}

// ----------------------------------------------------
// App: The main Onion application struct
// ----------------------------------------------------

// App represents the Onion application.
// It includes routing, middleware, and error handling capabilities.
type App struct {
	mux         *http.ServeMux // The HTTP multiplexer for routing
	middlewares []HandlerFunc  // Global middlewares applied to all requests
	notFound    HandlerFunc    // Custom 404 handler
}

// Route defines a single HTTP route.
// Method: HTTP method (e.g., GET, POST).
// Pattern: URL pattern (e.g., "/users/:id").
// Handler: Function to handle the route.
type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
}

// New creates a new Onion application with default settings.
// Initializes an empty router and a default 404 handler.
func New() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []HandlerFunc{},
		notFound: func(c *Context) {
			http.NotFound(c.Response, c.Request)
		},
	}
}

// Use adds a global middleware to the application.
// Middlewares are executed in the order they are added.
func (a *App) Use(mw HandlerFunc) {
	a.middlewares = append(a.middlewares, mw)
}

// NotFoundHandler sets a custom 404 handler for unknown routes.
func (a *App) NotFoundHandler(fn HandlerFunc) {
	a.notFound = fn
}

// MapRoutes maps multiple route slices to the application.
// Allows batch registration of routes, useful for grouping routes by feature.
func (a *App) MapRoutes(routeGroups ...[]Route) {
	for _, groupRoutes := range routeGroups {
		for _, r := range groupRoutes {
			a.handle(r.Method, r.Pattern, r.Handler)
		}
	}
}

// handle registers a single route with the application.
// Extracts path parameters and executes middlewares before calling the route handler.
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

// Run starts the Onion HTTP server on the specified address.
// Example: app.Run(":8080") to start on port 8080.
func (a *App) Run(addr string) error {
	fmt.Printf("Onion server running on %s\n", addr)
	return http.ListenAndServe(addr, a.mux)
}

// ----------------------------------------------------
// Helper Functions
// ----------------------------------------------------

// extractParams extracts path parameters from the URL based on the pattern.
// Example: For pattern "/users/:id" and path "/users/123", it returns {"id": "123"}.
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
