package books

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var res commonResponse
	var b bookReqValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res = commonResponse{"fail", err.Error(), nil}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := b.Validate()

	if !ok {
		res = commonResponse{"fail", msg, nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nb := bookModel{}
	mapBook(&nb, b)

	nb = saveBook(nb)
	bi := bookIdResponse{nb.Id}
	res = commonResponse{"success", "Book successfully added", bi.Response()}

	common.ResJSON(w, http.StatusCreated, res.Response())
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := commonResponse{"fail", err.Error(), nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	bq := GetBookQuery{q("name"), q("reading"), q("finished")}
	b := GetBooks(bq)
	bs := BooksSerializer{b}
	res := commonResponse{"success", "", bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	var res commonResponse
	p := common.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res = commonResponse{"fail", "Not Found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	b, ok := getBook(id)

	if !ok {
		res = commonResponse{"fail", "Not Found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{b}
	res = commonResponse{"success", "", bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func Put(w http.ResponseWriter, r *http.Request) {
	var res commonResponse
	var b bookReqValidator
	err := common.DecodeJSONBody(w, r, &b)
	if err != nil {
		res = commonResponse{"fail", err.Error(), nil}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := b.Validate()

	if !ok {
		res = commonResponse{"fail", msg, nil}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	p := common.ReqParams(r.URL.String())
	id := p("id")

	if id == "" {
		res = commonResponse{"fail", "Not Found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	nb := bookModel{}
	mapBook(&nb, b)

	nb, ok = updateBook(id, nb)

	if !ok {
		res = commonResponse{"fail", "Book with requested ID not found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{nb}
	res = commonResponse{"success", "Book successfully updated", bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var res commonResponse
	p := common.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res = commonResponse{"fail", "Not Found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	ob, ok := deleteBook(id)

	if !ok {
		res = commonResponse{"fail", "Book with requested ID not found", nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := BookSerializer{ob}
	res = commonResponse{"success", "Book successfully deleted", bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}
