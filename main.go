package main

import (
	"log"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
)

func main() {
	common.InitDB()

	common.RegisterHandler("/", "GET", common.RootPage)
	common.RegisterHandler("/books", "POST", books.Post)
	common.RegisterHandler("/books", "GET", books.Get)
	common.RegisterHandler("/books", "PUT", books.Put)
	common.RegisterHandler("/books", "DELETE", books.Delete)

	for v := range common.Routes {
		http.HandleFunc(v, common.MakeHandler)
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
