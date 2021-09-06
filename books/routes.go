package books

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var b bookReq
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: nil}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := b.Validate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

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
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	bq := GetBookQuery{q("name"), q("reading"), q("finished")}
	b := GetBooks(bq)
	bs := booksSerializer{b}
	res := common.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

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

	b, err := getBook(id)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{b}
	res := common.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func Put(w http.ResponseWriter, r *http.Request) {
	var b bookReq
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: nil}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := b.Validate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

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

	nb, err = updateBook(id, nb)

	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{nb}
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

	ob, err := deleteBook(id)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{ob}
	res := common.CommonResponse{Status: "success", Message: "Book successfully deleted", Data: bs.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}
