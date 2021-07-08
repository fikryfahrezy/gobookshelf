package main

import (
	"os"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
)

type Test struct {
	m string
}

func (t *Test) Write(p []byte) (n int, err error) {
	fo, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := fo.Write(p); err != nil {
		panic(err)
	}
	return 0, nil
}

func main() {
	common.InitDB()
	common.RegisterHandler("/", "GET", common.RootPage)
	common.RegisterHandler("/books", "POST", books.Post)
	common.RegisterHandler("/books", "GET", books.GetAll)
	common.RegisterHandler("/books/", "GET", books.GetOne)
	common.RegisterHandler("/books/", "PUT", books.Put)
	common.RegisterHandler("/books/", "DELETE", books.Delete)
	common.InitServer(8080)
}
