package main

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/pages"
)

func rootPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages", http.StatusFound)
}

func main() {
	common.InitDB()
	common.HanlderGET("/", rootPage)
	common.HanlderGET("/pages", pages.Page)
	common.HanlderGET("/matrix", pages.Matrix)
	common.HanlderPOST("/books", books.Post)
	common.HanlderGET("/books", books.GetAll)
	common.HanlderGET("/books/:id", books.GetOne)
	common.HanlderPUT("/books/:id", books.Put)
	common.HanlderDELETE("/books/:id", books.Delete)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	common.InitServer(8080)
}
