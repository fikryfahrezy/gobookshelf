package books

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"
)

// Testing Your (HTTP) Handlers in Go
// https://www.cloudbees.com/blog/testing-http-handlers-go
func TestGetRouteActive(t *testing.T) {
	common.Filename = "../data/books.json"

	common.RegisterHandler("/books", "POST", Post)
	// common.RegisterHandler("/books", "GET", Get)
	// common.RegisterHandler("/books", "PUT", Put)
	// common.RegisterHandler("/books", "DELETE", Delete)

	for u := range common.Routes {
		for m := range common.Routes[u] {
			// Use strings.NewReader() because:
			// https://golang.org/pkg/strings/#NewReader
			req, err := http.NewRequest(m, u, strings.NewReader(`{"name":"name","year":7,"author":"author","summary":"summary","publisher":"publisher","pageCount":0,"readPage":0,"reading":false}`))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			common.MakeHandler(rr, req)

			if rr.Result().StatusCode > 400 {
				t.FailNow()
			}
		}
	}
}
