package books

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/data"
)

func createBook(id string) bookModel {
	b := bookModel{
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
		"time",
		"time",
	}

	data.Insert(b)

	return b
}

// Testing Your (HTTP) Handlers in Go
// https://www.cloudbees.com/blog/testing-http-handlers-go
func TestHandlers(t *testing.T) {
	data.Filename = "../data/books.json"
	data.InitDB()

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
				createBook("5")
			},
			"/books?name=Name",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"Success, Book With Required ID found",
			func() {
				createBook("1")
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
				createBook("2")
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
				createBook("3")
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
			"PATCH",
			``,
			http.StatusMethodNotAllowed,
		},
	}

	common.HandlerPOST("/books", Post)
	common.HandlerGET("/books", GetAll)
	common.HandlerGET("/books/:id", GetOne)
	common.HandlerPUT("/books/:id", Put)
	common.HandlerDELETE("/books/:id", Delete)

	for _, c := range cases {
		c.init()

		// Use strings.NewReader() because:
		// https://golang.org/pkg/strings/#NewReader
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		common.MakeHandler(rr, req)

		// For debugging purpose
		// resp := rr.Result()
		// body, _ := io.ReadAll(resp.Body)
		// t.Log(resp.StatusCode)
		// t.Log(resp.Header.Get("Content-Type"))
		// t.Log(string(body))

		if rr.Result().StatusCode != c.expectedCode {
			t.FailNow()
		}
	}
}
