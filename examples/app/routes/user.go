package routes

import (
	"net/http"

	onion "github.com/saeedalam/Onion"
)

func GetAllUsers(c *onion.Context) {
	c.String(http.StatusOK, "Returning all users")
}

func GetSingleUser(c *onion.Context) {
	userID := c.Param("userId")
	c.String(http.StatusOK, "User ID: "+userID)
}

var UserRoutes = []onion.Route{
	{
		Method:  "GET",
		Pattern: "/users",
		Handler: GetAllUsers,
	},
	{
		Method:  "GET",
		Pattern: "/users/:userId",
		Handler: GetSingleUser,
	},
}
