package geocodings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func GetCountries(w http.ResponseWriter, r *http.Request) {
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
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
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	c, ce := make(chan []byte, 1), make(chan bool, 1)

	go func(cb chan []byte, cr chan bool) {
		resp, err := client.Do(req)
		if err != nil {
			cb <- []byte("")
			cr <- false
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			cb <- []byte("")
			cr <- false
			return
		}

		defer resp.Body.Close()

		cb <- body
		cr <- true
	}(c, ce)

	rs, ok := <-c, <-ce

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "something wrong", Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	var res interface{}
	json.Unmarshal(rs, &res)

	common.ResJSON(w, http.StatusOK, res)
}

func GetStreet(w http.ResponseWriter, r *http.Request) {
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	client := &http.Client{}
	rg, s := q("region"), q("street")

	if rg == "" || s == "" {
		res := common.CommonResponse{Status: "fail", Message: "query needed", Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://geocode.xyz/?geoit=json&region=%s&streetname=%s", rg, s), nil)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	defer resp.Body.Close()

	var res interface{}
	json.Unmarshal(body, &res)

	common.ResJSON(w, http.StatusOK, res)
}
