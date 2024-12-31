package routes

import (
	"net/http"
)

func GetAllBooks(c *onion.Context) {
	c.String(http.StatusOK, "Returning all books")
}

func GetBook(c *onion.Context) {
	bookID := c.Param("bookId")
	c.String(http.StatusOK, "Book ID: "+bookID)
}

var BookRoutes = []onion.Route{
	{
		Method:  "GET",
		Pattern: "/books",
		Handler: GetAllBooks,
	},
	{
		Method:  "GET",
		Pattern: "/books/:bookId",
		Handler: GetBook,
	},
}
