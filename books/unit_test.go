package books_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gosrouter"

	books_app "github.com/fikryfahrezy/gobookshelf/books/application"
	"github.com/fikryfahrezy/gobookshelf/books/domain/books"
	books_infra "github.com/fikryfahrezy/gobookshelf/books/infrastructure/books"
	books_http "github.com/fikryfahrezy/gobookshelf/books/interfaces/http"
)

func createBook(id string) books.BookModel {
	b := books.BookModel{
		id,
		"Name",
		1234,
		"Author",
		"Summary",
		"Publisher",
		1,
		1,
		false,
		false,
		false,
		"time",
		"time",
	}
	return b
}

// Testing Your (HTTP) Handlers in Go
// https://www.cloudbees.com/blog/testing-http-handlers-go
func TestBooks(t *testing.T) {
	bi := books_infra.InitDB("../data/books.json")
	ba := books_app.BookService{Fr: bi}
	bh := books_http.BookResource{Service: ba}
	books_http.AddRoutes(bh)

	cases := []struct {
		testName              string
		init                  func()
		url, method, bodydata string
		expectedCode          int
	}{
		{
			"Post Book Success",
			func() {},
			"/books",
			"POST",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusCreated,
		},
		{
			"Post Book Fail, Not Provided Body Data",
			func() {},
			"/books",
			"POST",
			``,
			http.StatusBadRequest,
		},
		{
			"Success, Without Query",
			func() {},
			"/books",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, With Query `?reading=0`",
			func() {},
			"/books?reading=0",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, With Query `?reading=1`",
			func() {},
			"/books?reading=1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, With Query `?finished=0`",
			func() {},
			"/books?finished=0",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, With Query `?finished=1`",
			func() {},
			"/books?finished=1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, With Query `?name=Name`",
			func() {
				bi.Save(createBook("5"))
			},
			"/books?name=Name",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, Book With Required ID found",
			func() {
				bi.Insert(createBook("1"))
			},
			"/books/1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Fail, Book With Required ID not found",
			func() {},
			"/books/10",
			"GET",
			``,
			http.StatusNotFound,
		},
		{
			"Fail, Route Not Found",
			func() {},
			"/books/",
			"GET",
			``,
			http.StatusNotFound,
		},
		{
			"Success, Book Updated",
			func() {
				bi.Insert(createBook("2"))
			},
			"/books/2",
			"PUT",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusOK,
		},
		{
			"Fail, Book Not Found",
			func() {},
			"/books/9",
			"PUT",
			``,
			http.StatusBadRequest,
		},
		{
			"Fail, Book Not Found",
			func() {},
			"/books/9",
			"PUT",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusNotFound,
		},
		{
			"Success, Book Deleted",
			func() {
				bi.Insert(createBook("3"))
			},
			"/books/3",
			"DELETE",
			``,
			http.StatusOK,
		},
		{
			"Success, Book Not Found",
			func() {},
			"/books/9",
			"DELETE",
			``,
			http.StatusNotFound,
		},
		{
			"Fail, Method Not Exist on Route",
			func() {},
			"/books",
			"NOAVA",
			``,
			http.StatusMethodNotAllowed,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.init()

			// Use strings.NewReader() because:
			// https://golang.org/pkg/strings/#NewReader
			req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			gosrouter.MakeHandler(rr, req)

			// For debugging purpose
			// resp := rr.Result()
			// body, _ := io.ReadAll(resp.Body)
			// t.Log(resp.StatusCode)
			// t.Log(resp.Header.Get("Content-Type"))
			// t.Log(string(body))

			if rr.Result().StatusCode != c.expectedCode {
				t.Fatal(rr.Result().StatusCode)
			}
		})
	}
}
