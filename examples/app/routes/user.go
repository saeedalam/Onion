package routes

import (
	"net/http"
	"onion"
)

var UserRoutes = onion.NewGroup("users").
	GET("/", GetAllUsers).
	GET("/:userId", GetUser).
	POST("/", CreateUser).
	PUT("/:userId", UpdateUser).
	DELETE("/:userId", DeleteUser).
	Routes()

func GetAllUsers(c *onion.Context) {
	c.String(http.StatusOK, "Returning all users")
}

func GetUser(c *onion.Context) {
	userID := c.Param("userId")
	c.String(http.StatusOK, "User ID: "+userID)
}

func CreateUser(c *onion.Context) {
	c.String(http.StatusOK, "Creating a new user")
}

func UpdateUser(c *onion.Context) {
	userID := c.Param("userId")
	c.String(http.StatusOK, "Updating user with ID: "+userID)
}

func DeleteUser(c *onion.Context) {
	userID := c.Param("userId")
	c.String(http.StatusOK, "Deleting user with ID: "+userID)
}
