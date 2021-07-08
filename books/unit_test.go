package books

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func createBook(id string) BookModel {
	b := BookModel{
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
	b.Save()

	return b
}

// Testing Your (HTTP) Handlers in Go
// https://www.cloudbees.com/blog/testing-http-handlers-go
func TestHandlers(t *testing.T) {
	common.Filename = "../data/books.json"
	cases := []struct {
		init                  func()
		url, method, bodydata string
		expectedCode          int
	}{
		{
			func() {},
			"/books",
			"POST",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusCreated,
		},
		{
			func() {},
			"/books",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {
				createBook("1")
			},
			"/books/1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {
				createBook("2")
			},
			"/books/2",
			"PUT",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusOK,
		},
		{
			func() {
				createBook("3")
			},
			"/books/3",
			"DELETE",
			``,
			http.StatusOK,
		},
		{
			func() {},
			"/books/",
			"GET",
			``,
			http.StatusNotFound,
		},
		{
			func() {},
			"/books",
			"PATCH",
			``,
			http.StatusMethodNotAllowed,
		},
		{
			func() {},
			"/books?reading=0",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {},
			"/books?reading=1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {},
			"/books?finished=0",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {},
			"/books?finished=1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			func() {
				createBook("5")
			},
			"/books?name=Name",
			"GET",
			``,
			http.StatusOK,
		},
	}

	common.RegisterHandler("/books", "POST", Post)
	common.RegisterHandler("/books", "GET", GetAll)
	common.RegisterHandler("/books/", "GET", GetOne)
	common.RegisterHandler("/books/", "PUT", Put)
	common.RegisterHandler("/books/", "DELETE", Delete)

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
		// fmt.Println(resp.StatusCode)
		// fmt.Println(resp.Header.Get("Content-Type"))
		// fmt.Println(string(body))

		if rr.Result().StatusCode != c.expectedCode {
			t.FailNow()
		}
	}
}
