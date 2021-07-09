package main

import (
	"io"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
)

func rootPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func main() {
	common.InitDB()
	common.HanlderGET("/", rootPage)
	common.HanlderPOST("/books", books.Post)
	common.HanlderGET("/books", books.GetAll)
	common.HanlderGET("/books/:id", books.GetOne)
	common.HanlderPUT("/books/:id", books.Put)
	common.HanlderDELETE("/books/:id", books.Delete)
	common.InitServer(8080)
}
