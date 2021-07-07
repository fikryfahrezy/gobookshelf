package books

import (
	"encoding/json"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var b BookModelValidator
	var res CommonResponse
	err := common.DecodeBody(w, r, &b)
	if err != nil {
		res = CommonResponse{"fail", err.Message, nil}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(res.Response())
		return
	}

	nb := saveBook(&b)
	bi := BookIdResponse{nb.Id}
	res = CommonResponse{"success", "Book successfully added", bi.Response()}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res.Response())
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b := getBooks()
	bs := BooksSerializer{b}
	res := CommonResponse{"success", "Book successfully added", bs.Response()}
	json.NewEncoder(w).Encode(res.Response())
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var res CommonResponse
	_, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Message, nil}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(res.Response())
		return
	}
	res = CommonResponse{"success", "Book successfully added", nil}
	json.NewEncoder(w).Encode(res.Response())
}

func Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var res CommonResponse
	_, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Message, nil}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(res.Response())
		return
	}
	res = CommonResponse{"success", "Book successfully updated", nil}
	json.NewEncoder(w).Encode(res.Response())
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var res CommonResponse
	_, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Message, nil}
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(res.Response())
		return
	}
	res = CommonResponse{"success", "Book successfully deleted", nil}
	json.NewEncoder(w).Encode(res.Response())
}
