package main

import (
	"example/middlewares"
	"example/routes"
)

func main() {
	app := onion.New()
	app.Use(middlewares.Auth)
	app.Use(middlewares.Log)

	app.MapRoutes(routes.UserRoutes, routes.BookRoutes)

	app.NotFoundHandler(func(c *onion.Context) {
		c.String(404, "Custom 404 message!")
	})

	app.Run(":8080")
}
