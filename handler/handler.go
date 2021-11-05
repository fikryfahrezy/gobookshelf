package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MalformedRequest struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}

type CommonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *CommonResponse) Response() *CommonResponse {
	return c
}

// How to Parse a JSON Request Body in Go (with Validation)
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) *MalformedRequest {
	defer r.Body.Close()

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

func AllowCORS(w *http.ResponseWriter) {
	rw := *w
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Request-Methods", "GET")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
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
