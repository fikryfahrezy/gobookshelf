package books

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var res CommonResponse
	var b BookModelValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res = CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, err.Status, res.Response())
		return
	}
	msg, ok := b.Validate()
	if !ok {
		res = CommonResponse{"fail", msg, nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}
	nb := saveBook(b)
	bi := BookIdResponse{nb.Id}
	res = CommonResponse{"success", "Book successfully added", bi.Response()}
	common.ResJSON(w, http.StatusCreated, res.Response())
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := getAllQuery(r.URL.String())
	if err != nil {
		res := CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}
	b := getBooks(q)
	bs := BooksSerializer{b}
	res := CommonResponse{"success", "", bs.Response()}
	common.ResJSON(w, http.StatusOK, res.Response())
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	var res CommonResponse
	id, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, err.Status, res.Response())
		return
	}
	b, ok := getBook(id)
	if !ok {
		res = CommonResponse{"fail", "Not Found", nil}
		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}
	bs := BookSerializer{b}
	res = CommonResponse{"success", "", bs.Response()}
	common.ResJSON(w, http.StatusOK, res.Response())
}

func Put(w http.ResponseWriter, r *http.Request) {
	var res CommonResponse
	var b BookModelValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res = CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, err.Status, res.Response())
		return
	}
	msg, ok := b.Validate()
	if !ok {
		res = CommonResponse{"fail", msg, nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}
	id, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, err.Status, res.Response())
		return
	}
	nb, ok := updateBook(id, b)
	if !ok {
		res = CommonResponse{"fail", "Book with requested ID not found", nil}
		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}
	bs := BookSerializer{nb}
	res = CommonResponse{"success", "Book successfully updated", bs.Response()}
	common.ResJSON(w, http.StatusOK, res.Response())
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var res CommonResponse
	id, err := common.RouteId(w, r.URL.Path)
	if err != nil {
		res = CommonResponse{"fail", err.Error(), nil}
		common.ResJSON(w, err.Status, res.Response())
		return
	}
	ob, ok := deleteBook(id)
	if !ok {
		res = CommonResponse{"fail", "Book with requested ID not found", nil}
		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}
	bs := BookSerializer{ob}
	res = CommonResponse{"success", "Book successfully deleted", bs.Response()}
	common.ResJSON(w, http.StatusOK, res.Response())
}
