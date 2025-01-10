
# Onion Wrapper

Onion is a lightweight, minimalist wrapper around Go’s built-in `net/http`. It aims to simplify the creation of web servers with an Express.js-like API and minimal boilerplate.

## Features
- **Easy Routing**: Define routes with simple chaining (`NewGroup("prefix").GET(...)`).  
- **Middleware**: Global middleware (`app.Use(...)`) or route groups.  
- **JSON & String Helpers**: Quickly serialize JSON or return plain text.  
- **Path Parameters**: Extract parameters like `/:id` into `c.Param("id")`.  
- **Custom 404**: Override the default “not found” behavior.

## Installation

```bash
go get github.com/saeedalam/Onion
```

## Quick Example

### 1. Project Layout

```
myproject/
  ├─ main.go
  ├─ routes/
  │   └─ books.go
  └─ middlewares/
      └─ auth.go
```

### 2. Define Your Routes

In `routes/books.go`:

```go
package routes

import (
    "net/http"
    "github.com/saeedalam/Onion"
)

// Handlers
func GetAllBooks(c *onion.Context) {
    c.String(http.StatusOK, "All books!")
}

func GetBook(c *onion.Context) {
    bookID := c.Param("bookId")
    c.String(http.StatusOK, "Book ID: " + bookID)
}

func CreateBook(c *onion.Context) {
    c.String(http.StatusOK, "Creating a new book!")
}

func UpdateBook(c *onion.Context) {
    bookID := c.Param("bookId")
    c.String(http.StatusOK, "Updating book with ID: " + bookID)
}

func DeleteBook(c *onion.Context) {
    bookID := c.Param("bookId")
    c.String(http.StatusOK, "Deleting book with ID: " + bookID)
}

// RouteGroup usage for prefix "books"
var BookRoutes = onion.NewGroup("books").
    GET("",          GetAllBooks).       // => GET /books
    GET("/:bookId",  GetBook).           // => GET /books/:bookId
    POST("",         CreateBook).        // => POST /books
    PUT("/:bookId",  UpdateBook).        // => PUT /books/:bookId
    DELETE("/:bookId", DeleteBook).      // => DELETE /books/:bookId
    Routes()
```

> **Note**: We used `.GET("")` rather than `.GET("/")` to avoid a trailing slash in the path.

### 3. (Optional) Middleware

In `middlewares/auth.go`:

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
    // If valid, execution continues to next middleware or handler
}
```

### 4. Main

In `main.go`:

```go
package main

import (
    "fmt"
    "github.com/saeedalam/Onion"
    "myproject/middlewares"
    "myproject/routes"
)

func main() {
    app := onion.New()

    // Global middleware
    app.Use(func(c *onion.Context) {
        fmt.Println("Executing global middleware")
    })
    app.Use(middlewares.Auth)

    // Register route groups
    app.UseRoutes(
        routes.BookRoutes,
        // Add more, e.g. routes.UserRoutes
    )

    // Custom 404
    app.NotFoundHandler(func(c *onion.Context) {
        c.String(http.StatusNotFound, "Custom 404: route not found!")
    })

    // Start server
    app.Run(":3333")
}
```

Now run:

```bash
go run main.go
```

**Test** with:

```bash
curl http://localhost:3333/books
curl http://localhost:3333/books/123
curl -X POST http://localhost:3333/books
curl -X PUT http://localhost:3333/books/123
curl -X DELETE http://localhost:3333/books/123
```

---

## Highlights

- **`RouteGroup`** allows prefix-based route definitions: `NewGroup("books").GET(...)`.  
- **`UseRoutes(...)`** bulk-registers route slices in one call.  
- **Global middleware** is applied in the order you call `app.Use(...)`.  
- **Path params** like `/:bookId` become `c.Param("bookId")`.  

That’s it! Enjoy a simpler, minimal “Onion” server for your Go web apps.
```
