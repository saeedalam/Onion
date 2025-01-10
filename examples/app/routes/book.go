package routes

import (
	"net/http"
	"onion" // your local module path
)

var BookRoutes = onion.NewGroup("books").
	GET("/", GetAllBooks).
	GET("/:bookId", GetBook).
	POST("/", CreateBook).
	PUT("/:bookId", UpdateBook).
	DELETE("/:bookId", DeleteBook).
	Routes()

func GetAllBooks(c *onion.Context) {
	c.String(http.StatusOK, "GET /books -> returning all books")
}

func GetBook(c *onion.Context) {
	id := c.Param("bookId")
	c.String(http.StatusOK, "GET /books/"+id+" -> single book")
}

func CreateBook(c *onion.Context) {
	c.String(http.StatusOK, "POST /books -> creating a book")
}

func UpdateBook(c *onion.Context) {
	id := c.Param("bookId")
	c.String(http.StatusOK, "PUT /books/"+id+" -> updating a book")
}

func DeleteBook(c *onion.Context) {
	id := c.Param("bookId")
	c.String(http.StatusOK, "DELETE /books/"+id+" -> deleting a book")
}
