package handler

import (
	"net/http"
	"reflect"
	"testing"
)

func getOne(w http.ResponseWriter, r *http.Request) {}

func getTwo(w http.ResponseWriter, r *http.Request) {}

func getThree(w http.ResponseWriter, r *http.Request) {}

func getFour(w http.ResponseWriter, r *http.Request) {}

func postOne(w http.ResponseWriter, r *http.Request) {}

func postTwo(w http.ResponseWriter, r *http.Request) {}

func postThree(w http.ResponseWriter, r *http.Request) {}

func postFour(w http.ResponseWriter, r *http.Request) {}

func putOne(w http.ResponseWriter, r *http.Request) {}

func putTwo(w http.ResponseWriter, r *http.Request) {}

func deleteOne(w http.ResponseWriter, r *http.Request) {}

func deleteTwo(w http.ResponseWriter, r *http.Request) {}

func TestGetRoute(t *testing.T) {
	Routes = make(map[string]RouteChild)

	cases := []struct {
		regUrl, reqUrl, mtd string
		regFn               func(url string, fn http.HandlerFunc)
		fn                  func(http.ResponseWriter, *http.Request)
	}{
		{
			"/",
			"/",
			"POST",
			HandlerPOST,
			postOne,
		},
		{
			"/",
			"/",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"/:id",
			"/1",
			"POST",
			HandlerPOST,
			postTwo,
		},
		{
			"/:id",
			"/1",
			"GET",
			HandlerGET,
			getTwo,
		},
		{
			"/:id",
			"/1",
			"PUT",
			HandlerPUT,
			putOne,
		},
		{
			"/:id",
			"/1",
			"DELETE",
			HandlerDELETE,
			deleteOne,
		},
		{
			"/one",
			"/one",
			"POST",
			HandlerPOST,
			postThree,
		},
		{
			"/one",
			"/one",
			"GET",
			HandlerGET,
			getThree,
		},
		{
			"/one/:id",
			"/one/1",
			"POST",
			HandlerPOST,
			postFour,
		},
		{
			"/one/:id",
			"/one/1",
			"GET",
			HandlerGET,
			getFour,
		},
		{
			"/one/:id",
			"/one/1",
			"PUT",
			HandlerPUT,
			putTwo,
		},
		{
			"/one/:id",
			"/one/1",
			"DELETE",
			HandlerDELETE,
			deleteTwo,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		if rt := getRoute(v.reqUrl, v.mtd); reflect.ValueOf(rt).Pointer() != reflect.ValueOf(v.fn).Pointer() {
			t.FailNow()
		}
	}
}

func TestDynamicRoute(t *testing.T) {
	Routes = make(map[string]RouteChild)

	cases := []struct {
		regUrl    string
		reqUrl    string
		paramName string
		param     string
		mtd       string
		regFn     func(url string, fn http.HandlerFunc)
		fn        func(http.ResponseWriter, *http.Request)
	}{
		{
			"/:id",
			"/1",
			"id",
			"1",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"/v1/:x",
			"/v1/11",
			"x",
			"11",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"/v2/v3/:xy",
			"/v2/v3/xyz",
			"xy",
			"xyz",
			"GET",
			HandlerGET,
			getOne,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		p := ReqParams(v.reqUrl)

		if p(v.paramName) != v.param {
			t.FailNow()
		}
	}
}

func TestReqQuery(t *testing.T) {
	Routes = make(map[string]RouteChild)

	cases := []struct {
		regUrl    string
		reqUrl    string
		queryName string
		query     string
		mtd       string
		regFn     func(url string, fn http.HandlerFunc)
		fn        func(http.ResponseWriter, *http.Request)
	}{
		{
			`/`,
			`/?q=a`,
			"q",
			"a",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"/v1",
			`/v1?q=a&b=b`,
			"b",
			"b",
			"GET",
			HandlerGET,
			getOne,
		},
		{
			"/v2/v3",
			`/v2/v3/xyz?a=a&b=b&c=c`,
			"c",
			"c",
			"GET",
			HandlerGET,
			getOne,
		},
	}

	for _, v := range cases {
		v.regFn(v.regUrl, v.fn)
	}

	for _, v := range cases {
		p, _ := ReqQuery(v.reqUrl)

		if p(v.queryName) != v.query {
			t.FailNow()
		}
	}
}
