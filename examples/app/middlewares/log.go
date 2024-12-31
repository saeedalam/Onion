package middlewares

import (
    "fmt"
    "github.com/saeedalam/Onion"
)

func Log(c *onion.Context) {
    fmt.Printf("[Log Middleware] %s %s
", c.Request.Method, c.Request.URL.Path)
}
