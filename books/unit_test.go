package books

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"
)

// Testing Your (HTTP) Handlers in Go
// https://www.cloudbees.com/blog/testing-http-handlers-go
func TestHandlers(t *testing.T) {
	common.Filename = "../data/books.json"
	cases := []struct {
		url, method, bodydata string
		expectedCode          int
	}{
		{
			"/books",
			"POST",
			`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`,
			http.StatusCreated,
		},
		{
			"/books",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"/books/1",
			"GET",
			``,
			http.StatusOK,
		},
		{
			"/books/1",
			"PUT",
			``,
			http.StatusOK,
		},
		{
			"/books/1",
			"DELETE",
			``,
			http.StatusOK,
		},
		{
			"/books/",
			"GET",
			``,
			http.StatusNotFound,
		},
		{
			"/books",
			"PATCH",
			``,
			http.StatusMethodNotAllowed,
		},
	}

	common.RegisterHandler("/books", "POST", Post)
	common.RegisterHandler("/books", "GET", GetAll)
	common.RegisterHandler("/books/", "GET", GetOne)
	common.RegisterHandler("/books/", "PUT", Put)
	common.RegisterHandler("/books/", "DELETE", Delete)

	for _, c := range cases {
		// Use strings.NewReader() because:
		// https://golang.org/pkg/strings/#NewReader
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		common.MakeHandler(rr, req)

		resp := rr.Result()
		body, _ := io.ReadAll(resp.Body)

		t.Log(resp.StatusCode)
		t.Log(resp.Header.Get("Content-Type"))
		t.Log(string(body))

		if rr.Result().StatusCode != c.expectedCode {
			t.FailNow()
		}
	}
}
