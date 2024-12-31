package middlewares

import (
	"fmt"
)

func Log(c *onion.Context) {
	fmt.Printf("[Log Middleware] %s %s", c.Request.Method, c.Request.URL.Path)
}
