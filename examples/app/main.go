package main

import (
    "example/routes"
    "example/middlewares"
    "github.com/saeedalam/Onion"
)

func main() {
    app := onion.New()
    app.Use(middlewares.Auth)
    app.Use(middlewares.Log)

    adminGroup := app.Group("/admin")
    adminGroup.Use(middlewares.AdminAuth)
    adminGroup.MapRoutes(routes.AdminRoutes)

    app.MapRoutes(routes.UserRoutes, routes.BookRoutes)

    app.NotFoundHandler(func(c *onion.Context) {
        c.String(404, "Custom 404 message!")
    })

    app.Run(":8080")
}
