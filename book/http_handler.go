package book

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"

	"github.com/fikryfahrezy/gosrouter"
)

type bookReq struct {
	Name      string `json:"name"`
	Year      int    `json:"year"`
	Author    string `json:"author"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
	PageCount int    `json:"pageCount"`
	ReadPage  int    `json:"readPage"`
	Reading   bool   `json:"reading"`
}

type bookResponse struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Year       int    `json:"year"`
	Author     string `json:"author"`
	Summary    string `json:"summary"`
	Publisher  string `json:"publisher"`
	PageCount  int    `json:"pageCount"`
	ReadPage   int    `json:"readPage"`
	Finished   bool   `json:"finished"`
	Reading    bool   `json:"reading"`
	IsDeleted  bool   `json:"isDeleted"`
	InsertedAt string `json:"insertedAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type bookIdResponse struct {
	BookId string `json:"bookId"`
}

type book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type booksResponse struct {
	Books []book `json:"books"`
}

func mapBookQueryToResponse(ob ResCommand) book {
	nb := book{
		Id:        ob.Id,
		Name:      ob.Name,
		Publisher: ob.Publisher,
	}
	return nb
}

func mapBookQueryToResponses(ob []ResCommand) []book {
	bs := make([]book, len(ob))
	for i, b := range ob {
		bs[i] = mapBookQueryToResponse(b)
	}

	return bs
}

type bookSerializer struct {
	Book bookResponse `json:"book"`
}

func mapBookReqToCmd(nb bookReq) ReqCommand {
	ob := ReqCommand{
		Name:      nb.Name,
		Year:      nb.Year,
		Author:    nb.Author,
		Summary:   nb.Summary,
		Publisher: nb.Publisher,
		PageCount: nb.PageCount,
		ReadPage:  nb.ReadPage,
		Finished:  nb.ReadPage == nb.PageCount,
		Reading:   nb.Reading,
	}
	return ob
}

func mapBookCmdToResponse(ob ResCommand) bookResponse {
	nb := bookResponse{
		Id:         ob.Id,
		Name:       ob.Name,
		Year:       ob.Year,
		Author:     ob.Author,
		Summary:    ob.Summary,
		Publisher:  ob.Publisher,
		PageCount:  ob.PageCount,
		ReadPage:   ob.ReadPage,
		Finished:   ob.Finished,
		Reading:    ob.Reading,
		IsDeleted:  ob.IsDeleted,
		InsertedAt: ob.InsertedAt,
		UpdatedAt:  ob.UpdatedAt,
	}
	return nb
}

func mapQueryToBookQueryCmd(n, r, f string) QueryCommand {
	q := QueryCommand{n, r, f}
	return q
}

func Post(w http.ResponseWriter, r *http.Request, br *Service) {
	var b bookReq
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: nil}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	nb, err := br.SaveBook(mapBookReqToCmd(b))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	bi := bookIdResponse{nb.Id}
	res := common.Response{Message: "Book successfully added", Data: bi}
	common.ResJSON(w, http.StatusCreated, res)
}

func GetAll(w http.ResponseWriter, r *http.Request, br *Service) {
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.Response{Message: err.Error(), Data: make([]interface{}, 0)}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	b := br.GetBooks(mapQueryToBookQueryCmd(q("name"), q("reading"), q("finished")))
	bs := make([]bookResponse, len(b))
	for i, bi := range b {
		bs[i] = mapBookCmdToResponse(bi)
	}

	bt := booksResponse{Books: mapBookQueryToResponses(b)}
	res := common.Response{Message: "", Data: bt}
	common.ResJSON(w, http.StatusOK, res)
}

func GetOne(w http.ResponseWriter, r *http.Request, br *Service) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := common.Response{Message: "Not Found", Data: nil}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	b, err := br.GetBook(id)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}

		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(b)}
	res := common.Response{Message: "", Data: bs}
	common.ResJSON(w, http.StatusOK, res)
}

func Put(w http.ResponseWriter, r *http.Request, br *Service) {
	var b bookReq
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: nil}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	p := gosrouter.ReqParams(r.URL.String())
	id := p("id")

	if id == "" {
		res := common.Response{Message: "Not Found", Data: nil}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	nb, err := br.UpdateBook(id, mapBookReqToCmd(b))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(nb)}
	res := common.Response{Message: "Book successfully updated", Data: bs}
	common.ResJSON(w, http.StatusOK, res)
}

func Delete(w http.ResponseWriter, r *http.Request, br *Service) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := common.Response{Message: "Not Found", Data: nil}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	ob, err := br.DeleteBook(id)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(ob)}
	res := common.Response{Message: "Book successfully deleted", Data: bs}
	common.ResJSON(w, http.StatusOK, res)
}

func AddRoutes(s *Service) {
	gosrouter.HandlerPOST("/books", s.Http(Post))
	gosrouter.HandlerGET("/books", s.Http(GetAll))
	gosrouter.HandlerGET("/books/:id", s.Http(GetOne))
	gosrouter.HandlerPUT("/books/:id", s.Http(Put))
	gosrouter.HandlerDELETE("/books/:id", s.Http(Delete))
}
