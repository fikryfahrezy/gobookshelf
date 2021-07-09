package common

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var routeMethods = map[string]string{
	"GET":    "GET",
	"POST":   "POST",
	"PUT":    "PUT",
	"DELETE": "DELETE",
}

type RouteChild struct {
	Depth int
	Route string
	Val   string
	Fn    map[string]func(http.ResponseWriter, *http.Request)
	Child *RouteChild
}

func (r *RouteChild) CreateFn(mtd string, fn func(http.ResponseWriter, *http.Request)) {
	if f := r.Fn; f == nil {
		r.Fn = make(map[string]func(http.ResponseWriter, *http.Request))
	}
	r.Fn[mtd] = fn
}

// var Routes = make(map[string]map[string]func(http.ResponseWriter, *http.Request))

// Deprecated: Not using this anymore
// func oldRegisterHandler(mtd string, url string, fn http.HandlerFunc) {
// 	m := routeMethods[strings.ToUpper(mtd)]
// 	if m == "" {
// 		return
// 	}
// 	if ro := Routes[url]; ro == nil {
// 		Routes[url] = make(map[string]func(http.ResponseWriter, *http.Request))
// 	}
// 	Routes[url][mtd] = fn
// }

var Routes = make(map[string]RouteChild)

func routeChild(r *RouteChild, i int, m int, mtd string, s []string, fn func(http.ResponseWriter, *http.Request)) *RouteChild {
	nr := r
	if r == nil {
		nr = &RouteChild{}
	}
	if i == m {
		nr.CreateFn(mtd, fn)
		return nr
	}
	if i >= m {
		nr.Route = "/" + s[0]
		nr.CreateFn(mtd, fn)
		return nr
	}
	nr.Depth = i
	nr.Route = "/" + s[i]
	i++
	if nc := routeChild(nr.Child, i, m, mtd, s, fn); nc != nil {
		nr.Child = nc
	}
	return nr
}

func registerHandler(mtd string, url string, fn func(http.ResponseWriter, *http.Request)) {
	s := strings.Split(url, "/")
	l := len(s)
	if l <= 1 {
		return
	}
	s = s[1:]
	r := "/" + s[0]
	if o := Routes[r]; o.Route != "" {
		Routes[r] = *routeChild(&o, 1, l-1, mtd, s, fn)
	} else {
		Routes[r] = RouteChild{0, r, "", map[string]func(http.ResponseWriter, *http.Request){
			mtd: fn,
		}, nil}
	}
}

func getRoute(url string, mtd string) func(http.ResponseWriter, *http.Request) {
	s := strings.Split(url, "/")
	s = s[1:]
	var l RouteChild
	if len(s) == 1 {
		if h := Routes["/"+s[0]].Fn[mtd]; h != nil {
			return h
		}
	}
	for i, v := range s {
		v = "/" + v
		if i == 0 {
			l = Routes[v]
		}
		if f := l.Fn[mtd]; f != nil && i == len(s)-1 {
			return f
		}
		if l.Child != nil {
			l = *l.Child
		}
	}
	return nil
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	m := routeMethods[strings.ToUpper(r.Method)]
	if m == "" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Old Get Route
	// re := regexp.MustCompile(`^([^/]*/[^/]*/).*$`)
	// u := re.ReplaceAllString(r.URL.Path, "$1")
	// rt := Routes[u][m]
	rt := getRoute(r.URL.Path, m)
	if rt == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	rt(w, r)
}

func HanlderPOST(url string, fn http.HandlerFunc) {
	// oldRegisterHandler("POST", url, fn)
	registerHandler("POST", url, fn)
}

func HanlderGET(url string, fn http.HandlerFunc) {
	// oldRegisterHandler("GET", url, fn)
	registerHandler("GET", url, fn)
}

func HanlderPUT(url string, fn http.HandlerFunc) {
	// oldRegisterHandler("PUT", url, fn)
	registerHandler("PUT", url, fn)
}

func HanlderDELETE(url string, fn http.HandlerFunc) {
	// oldRegisterHandler("DELETE", url, fn)
	registerHandler("DELETE", url, fn)
}

func InitServer(p int) {
	for v := range Routes {
		http.HandleFunc(v, MakeHandler)
	}
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(p), nil))
}

func Test(url string, mtd string) func(http.ResponseWriter, *http.Request) {
	return getRoute(url, mtd)
}
