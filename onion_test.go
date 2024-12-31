package onion

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBasicRouting tests a basic route with a fixed path.
func TestBasicRouting(t *testing.T) {
	app := New()

	app.handle("GET", "/hello", func(c *Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()

	// Use app.mux.ServeHTTP directly
	app.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	if rec.Body.String() != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", rec.Body.String())
	}
}

// TestPathParameters tests a route with path parameters.
func TestPathParameters(t *testing.T) {
	app := New()

	app.handle("GET", "/users/:id", func(c *Context) {
		userID := c.Param("id")
		c.String(http.StatusOK, "User ID: "+userID)
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	rec := httptest.NewRecorder()

	// Use app.mux.ServeHTTP directly
	app.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	expected := "User ID: 123"
	if rec.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, rec.Body.String())
	}
}

// TestNotFound tests the custom 404 handler.
func TestNotFound(t *testing.T) {
	app := New()

	app.NotFoundHandler(func(c *Context) {
		c.String(http.StatusNotFound, "Custom 404")
	})

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()

	// Use app.mux.ServeHTTP directly
	app.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", rec.Code)
	}

	if rec.Body.String() != "Custom 404" {
		t.Errorf("Expected body 'Custom 404', got '%s'", rec.Body.String())
	}
}

// TestMiddleware tests global middleware functionality.
func TestMiddleware(t *testing.T) {
	app := New()

	// Add a simple middleware to add a header
	app.Use(func(c *Context) {
		c.Response.Header().Set("X-Test", "MiddlewarePassed")
	})

	app.handle("GET", "/middleware", func(c *Context) {
		c.String(http.StatusOK, "Middleware Test")
	})

	req := httptest.NewRequest("GET", "/middleware", nil)
	rec := httptest.NewRecorder()

	// Use app.mux.ServeHTTP directly
	app.mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	if rec.Header().Get("X-Test") != "MiddlewarePassed" {
		t.Errorf("Expected X-Test header to be 'MiddlewarePassed', got '%s'", rec.Header().Get("X-Test"))
	}
}

// TestMultipleRoutes ensures multiple routes can be registered and handled correctly.
func TestMultipleRoutes(t *testing.T) {
	app := New()

	app.handle("GET", "/route1", func(c *Context) {
		c.String(http.StatusOK, "Route 1")
	})
	app.handle("GET", "/route2", func(c *Context) {
		c.String(http.StatusOK, "Route 2")
	})

	req1 := httptest.NewRequest("GET", "/route1", nil)
	rec1 := httptest.NewRecorder()
	app.mux.ServeHTTP(rec1, req1)

	if rec1.Body.String() != "Route 1" {
		t.Errorf("Expected body 'Route 1', got '%s'", rec1.Body.String())
	}

	req2 := httptest.NewRequest("GET", "/route2", nil)
	rec2 := httptest.NewRecorder()
	app.mux.ServeHTTP(rec2, req2)

	if rec2.Body.String() != "Route 2" {
		t.Errorf("Expected body 'Route 2', got '%s'", rec2.Body.String())
	}
}
