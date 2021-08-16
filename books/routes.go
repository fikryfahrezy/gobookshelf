package books

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var b bookReqValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := b.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nb := bookModel{}
	mapBook(&nb, b)

	nb = saveBook(nb)
	bi := bookIdResponse{nb.Id}
	res := common.CommonResponse{Status: "success", Message: "Book successfully added", Data: bi.Response()}

	common.ResJSON(w, http.StatusCreated, res.Response())
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	bq := GetBookQuery{q("name"), q("reading"), q("finished")}
	b := GetBooks(bq)
	bs := BooksSerializer{b}
	res := common.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

	common.AllowCORS(&w)
	common.ResJSON(w, http.StatusOK, res.Response())
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	p := common.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := common.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	b, ok := getBook(id)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{b}
	res := common.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func Put(w http.ResponseWriter, r *http.Request) {
	var b bookReqValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := b.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	p := common.ReqParams(r.URL.String())
	id := p("id")

	if id == "" {
		res := common.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	nb := bookModel{}
	mapBook(&nb, b)

	nb, ok = updateBook(id, nb)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "Book with requested ID not found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{nb}
	res := common.CommonResponse{Status: "success", Message: "Book successfully updated", Data: bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func Delete(w http.ResponseWriter, r *http.Request) {
	p := common.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := common.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	ob, ok := deleteBook(id)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "Book with requested ID not found", Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{ob}
	res := common.CommonResponse{Status: "success", Message: "Book successfully deleted", Data: bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}
