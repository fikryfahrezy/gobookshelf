package http

import (
	"errors"
	"net/http"

	"github.com/fikryfahrezy/gosrouter"

	"github.com/fikryfahrezy/gobookshelf/books/application"

	"github.com/fikryfahrezy/gobookshelf/handler"
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

func (b *bookReq) Validate() error {
	if b.Name == "" {
		return errors.New("name cannot be empty")
	}

	if b.Author == "" {
		return errors.New("author cannot be empty")
	}

	if b.Summary == "" {
		return errors.New("summary cannot be empty")
	}

	if b.Publisher == "" {
		return errors.New("publisher cannot be empty")
	}

	if b.ReadPage > b.PageCount {
		return errors.New("read page cannot be bigger than page count")
	}

	return nil
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

func (b *bookIdResponse) Response() *bookIdResponse {
	return b
}

type book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type booksSerializer struct {
	Books []bookResponse
}

type booksResponse struct {
	Books []book `json:"books"`
}

func (s *booksSerializer) Response() booksResponse {
	b := booksResponse{
		Books: make([]book, len(s.Books)),
	}

	for i, v := range s.Books {
		b.Books[i] = book{v.Id, v.Name, v.Publisher}
	}

	return b
}

type bookSerializer struct {
	Book bookResponse `json:"book"`
}

func (b *bookSerializer) Response() bookSerializer {
	return *b
}

func mapBookReqToCmd(nb bookReq) application.BookReqCommand {
	ob := application.BookReqCommand{
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

func mapBookCmdToResponse(ob application.BookResCommand) bookResponse {
	nb := bookResponse{
		Id:         ob.Id,
		Name:       ob.Name,
		Year:       ob.Year,
		Author:     ob.Author,
		Summary:    ob.Summary,
		Publisher:  ob.Summary,
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

func mapQueryToBookQueryCmd(n, r, f string) application.BookQueryCommand {
	q := application.BookQueryCommand{n, r, f}
	return q
}

type BookResource struct {
	Service application.BookService
}

func (br BookResource) Post(w http.ResponseWriter, r *http.Request) {
	var b bookReq
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: nil}
		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := b.Validate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nb := br.Service.SaveBook(mapBookReqToCmd(b))
	bi := bookIdResponse{nb.Id}
	res := handler.CommonResponse{Message: "Book successfully added", Data: bi.Response()}
	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func (br BookResource) GetAll(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	q, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: make([]interface{}, 0)}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	b := br.Service.GetBooks(mapQueryToBookQueryCmd(q("name"), q("reading"), q("finished")))
	bs := make([]bookResponse, len(b))
	for i, bi := range b {
		bs[i] = mapBookCmdToResponse(bi)
	}

	bt := booksSerializer{bs}
	res := handler.CommonResponse{Message: "", Data: bt.Response()}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (br BookResource) GetOne(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := handler.CommonResponse{Message: "Not Found", Data: nil}
		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	b, err := br.Service.GetBook(id)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}

		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(b)}
	res := handler.CommonResponse{Message: "", Data: bs.Response()}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (br BookResource) Put(w http.ResponseWriter, r *http.Request) {
	var b bookReq
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: nil}
		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := b.Validate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	p := gosrouter.ReqParams(r.URL.String())
	id := p("id")

	if id == "" {
		res := handler.CommonResponse{Message: "Not Found", Data: nil}
		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	nb, err := br.Service.UpdateBook(id, mapBookReqToCmd(b))
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(nb)}
	res := handler.CommonResponse{Message: "Book successfully updated", Data: bs.Response()}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (br BookResource) Delete(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	id := p("id")

	if id == "" {
		res := handler.CommonResponse{Message: "Not Found", Data: nil}
		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	ob, err := br.Service.DeleteBook(id)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusNotFound, res.Response())
		return
	}

	bs := bookSerializer{Book: mapBookCmdToResponse(ob)}
	res := handler.CommonResponse{Message: "Book successfully deleted", Data: bs.Response()}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func AddRoutes(s BookResource) {
	gosrouter.HandlerPOST("/books", s.Post)
	gosrouter.HandlerGET("/books", s.GetAll)
	gosrouter.HandlerGET("/books/:id", s.GetOne)
	gosrouter.HandlerPUT("/books/:id", s.Put)
	gosrouter.HandlerDELETE("/books/:id", s.Delete)
}
