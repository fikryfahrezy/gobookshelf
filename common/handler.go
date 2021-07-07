package common

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

var routeMethods = map[string]string{
	"GET":    "GET",
	"POST":   "POST",
	"PUT":    "PUT",
	"DELETE": "DELETE",
}

var Routes = make(map[string]map[string]func(http.ResponseWriter, *http.Request))

func RootPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	m := routeMethods[strings.ToUpper(r.Method)]

	if m == "" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	re := regexp.MustCompile(`^([^/]*/[^/]*/).*$`)
	u := re.ReplaceAllString(r.URL.Path, "$1")
	rt := Routes[u][m]

	if rt == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	rt(w, r)
}

func RegisterHandler(url string, mtd string, fn http.HandlerFunc) {
	m := routeMethods[strings.ToUpper(mtd)]
	if m == "" {
		return
	}

	ro := Routes[url]
	if ro == nil {
		ro = make(map[string]func(http.ResponseWriter, *http.Request))
		Routes[url] = ro
	}

	Routes[url][mtd] = fn
}
