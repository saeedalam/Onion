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
