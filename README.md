# Onion Framework

Onion is a lightweight and minimalist wrapper for Go's `http` module. It aims to simplify the creation of web servers with an Express.js-like API and minimal boilerplate.

## Features
- **Easy Routing**: Define routes with simple `app.Get`, `app.Post`, etc.
- **Middleware**: Global or route-specific middleware with `app.Use`.
- **Route Groups**: Prefix-based route grouping with `app.Group`.
- **JSON Parsing**: Automatic decoding of JSON into structs.
- **Path Parameters**: Extract parameters like `/users/:id`.
- **NotFoundHandler**: Customizable 404 pages.

## Quick Start

### Install
```bash
go get github.com/saeedalam/Onion
```

### Example App

```go
package main

import (
    "github.com/saeedalam/Onion"
    "example/routes"
    "example/middlewares"
)

func main() {
    app := onion.New()

    // Add global middleware
    app.Use(middlewares.Auth)
    app.Use(middlewares.Log)

    // Define route groups
    admin := app.Group("/admin")
    admin.Use(middlewares.AdminAuth)
    admin.MapRoutes(routes.AdminRoutes)

    // Map general routes
    app.MapRoutes(routes.UserRoutes, routes.BookRoutes)

    // Set a custom 404 handler
    app.NotFoundHandler(func(c *onion.Context) {
        c.String(404, "Route not found!")
    })

    // Run the server
    app.Run(":8080")
}
```

### Define Routes
Create a `routes` folder with `user.go`:

```go
package routes

import (
    "net/http"
    "github.com/saeedalam/Onion"
)

func GetAllUsers(c *onion.Context) {
    c.String(http.StatusOK, "All users!")
}

func GetUser(c *onion.Context) {
    id := c.Param("id")
    c.String(http.StatusOK, "User ID: "+id)
}

var UserRoutes = []onion.Route{
    {Method: "GET", Pattern: "/users", Handler: GetAllUsers},
    {Method: "GET", Pattern: "/users/:id", Handler: GetUser},
}
```

### Middleware
Add middleware like `auth.go`:

```go
package middlewares

import (
    "net/http"
    "github.com/saeedalam/Onion"
)

func Auth(c *onion.Context) {
    token := c.Request.Header.Get("X-Auth")
    if token == "" {
        c.String(http.StatusUnauthorized, "Unauthorized!")
        return
    }
}
```
