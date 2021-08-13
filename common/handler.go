package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var routeMethods = map[string]string{
	"GET":    "GET",
	"POST":   "POST",
	"PUT":    "PUT",
	"DELETE": "DELETE",
}

type RouteChild struct {
	Depth   int
	Route   string
	Dynamic bool
	Fn      map[string]func(http.ResponseWriter, *http.Request)
	Child   *RouteChild
}

func (r *RouteChild) CreateFn(mtd string, fn func(http.ResponseWriter, *http.Request)) {
	if f := r.Fn; f == nil {
		r.Fn = make(map[string]func(http.ResponseWriter, *http.Request))
	}

	r.Fn[mtd] = fn
}

var Routes = make(map[string]RouteChild)

func routeChild(r *RouteChild, i, m int, mtd string, s []string, fn func(http.ResponseWriter, *http.Request)) *RouteChild {
	nr := r

	if r == nil {
		nr = &RouteChild{}
	}

	rt := s[i]
	nr.Depth = i
	nr.Route = "/" + rt

	if strings.HasPrefix(rt, ":") {
		nr.Dynamic = true
	} else {
		nr.Dynamic = false
	}

	if i == m-1 {
		nr.CreateFn(mtd, fn)
		return nr
	}

	i++

	if nc := routeChild(nr.Child, i, m, mtd, s, fn); nc != nil {
		nr.Child = nc
	}

	return nr
}

func registerHandler(mtd, url string, fn func(http.ResponseWriter, *http.Request)) {
	if strings.Contains(url, ":") {
		s := strings.Split(url, "/")
		s = s[1:]
		l := len(s)

		if l == 0 {
			return
		}

		fe := s[0]
		r := "/"

		if !strings.HasPrefix(fe, ":") {
			r += fe
		} else {
			// Ref: How to prepend int to slice
			// https://stackoverflow.com/questions/53737435/how-to-prepend-int-to-slice
			s = append(s, "")
			copy(s[1:], s)
			s[0] = ""
			l = len(s)
		}

		o := Routes[r]
		Routes[r] = *routeChild(&o, 0, l, mtd, s, fn)
	} else {
		o := Routes[url]

		if o.Route == "" {
			o = RouteChild{Depth: 0, Route: url}
		}

		o.CreateFn(mtd, fn)

		if Routes[url].Route == "" {
			Routes[url] = o
		}
	}
}

func HandlerPOST(url string, fn http.HandlerFunc) {
	registerHandler("POST", url, fn)
}

func HandlerGET(url string, fn http.HandlerFunc) {
	registerHandler("GET", url, fn)
}

func HandlerPUT(url string, fn http.HandlerFunc) {
	registerHandler("PUT", url, fn)
}

func HandlerDELETE(url string, fn http.HandlerFunc) {
	registerHandler("DELETE", url, fn)
}

func getRoute(url, mtd string) func(http.ResponseWriter, *http.Request) {
	if r := Routes[url]; r.Route == url && r.Fn[mtd] != nil {
		return r.Fn[mtd]
	}

	s := strings.Split(url, "/")
	s = s[1:]

	if len(s) == 1 {
		if h := Routes["/"+s[0]].Fn[mtd]; h != nil {
			return h
		}
	}

	var l RouteChild

	for i, v := range s {
		v = "/" + v

		if i == 0 {
			if r, rc := Routes[v], Routes["/"].Child; r.Route == "" && rc != nil {
				l = *rc
			} else {
				l = r
			}
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

	rt := getRoute(r.URL.Path, m)

	if rt == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	rt(w, r)
}

func InitServer(p int) {
	for v := range Routes {
		http.HandleFunc(v, MakeHandler)
	}

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(p), nil))
}

type MalformedRequest struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}

// How to Parse a JSON Request Body in Go (with Validation)
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *MalformedRequest {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	cntType := r.Header.Get("Content-Type")

	if cntType != "" {
		if cntType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &MalformedRequest{http.StatusUnsupportedMediaType, msg}
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

			return &MalformedRequest{http.StatusBadRequest, msg}

		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"

			return &MalformedRequest{http.StatusBadRequest, msg}

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in Requeired struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)

			return &MalformedRequest{http.StatusBadRequest, msg}

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in Required struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)

			return &MalformedRequest{http.StatusBadRequest, msg}

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"

			return &MalformedRequest{http.StatusBadRequest, msg}

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"

			return &MalformedRequest{http.StatusRequestEntityTooLarge, msg}

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			return &MalformedRequest{http.StatusInternalServerError, err.Error()}
		}
	}

	// Otherwise default to logging the error and sending a 500 Internal
	// Server Error response.
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		err = &MalformedRequest{http.StatusBadRequest, msg}
	}

	return nil
}

func ReqQuery(r string) (func(n string) string, error) {
	u, err := url.Parse(r)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	return func(n string) string {
		return q.Get(n)
	}, nil
}

func ReqParams(u string) func(p string) string {
	s := strings.Split(u, "/")
	s = s[1:]

	return func(p string) string {
		var l RouteChild
		isSls := false

		for i, v := range s {
			v = "/" + v

			if i == 0 {
				if r := Routes[v]; r.Route == "" {
					l = *Routes["/"].Child
					isSls = true
				} else {
					l = r
				}
			}

			if l.Dynamic && strings.Split(l.Route, "/:")[1] == p {
				if isSls {
					return s[l.Depth-1]
				}
				return s[l.Depth]
			}

			if l.Child != nil {
				l = *l.Child
			}
		}

		return ""
	}
}

func ResJSON(w http.ResponseWriter, s int, v interface{}) {
	c := w.Header().Get("Content-Type")

	if c == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	} else if !strings.Contains(c, "application/json") {
		w.Header().Add("Content-Type", "application/json")
	}

	w.WriteHeader(s)

	json.NewEncoder(w).Encode(v)
}
