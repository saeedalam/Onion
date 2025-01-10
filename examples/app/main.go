package main

import (
	"fmt"
	"onion"
	"onion/examples/app/routes"
)

func main() {
	app := onion.New()

	// Global middleware
	app.Use(func(c *onion.Context) {
		fmt.Println("Executing global middleware")
	})

	// Register routes
	app.UseRoutes(
		routes.BookRoutes,
		routes.UserRoutes,
	)

	// Start the server
	app.Run(":3333")
}
