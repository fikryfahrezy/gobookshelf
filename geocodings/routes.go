package geocodings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/handler"
)

func GetCountries(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	q, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: make([]interface{}, 0)}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	var req *http.Request
	client := &http.Client{}
	n := q("name")

	if n == "" {
		req, err = http.NewRequest("GET", "https://restcountries.com/v2/all", nil)
	} else {
		req, err = http.NewRequest("GET", fmt.Sprintf("https://restcountries.com/v2/name/%s", n), nil)
	}

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: make([]interface{}, 0)}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	var res interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	defer resp.Body.Close()

	handler.ResJSON(w, http.StatusOK, res)
}

func GetStreet(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	q, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	client := &http.Client{}
	rg, s := q("region"), q("street")

	if rg == "" || s == "" {
		res := handler.CommonResponse{Message: "query needed", Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://geocode.xyz/?geoit=json&region=%s&streetname=%s", rg, s), nil)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	var res interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	defer resp.Body.Close()

	handler.ResJSON(w, http.StatusOK, res)
}
