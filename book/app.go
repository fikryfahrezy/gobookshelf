package book

import (
	"errors"
	"fmt"
	"net/http"
)

type QueryCommand struct {
	Name     string
	Reading  string
	Finished string
}

type ReqCommand struct {
	Id        string
	Name      string
	Year      int
	Author    string
	Summary   string
	Publisher string
	PageCount int
	ReadPage  int
	Finished  bool
	Reading   bool
}

func (b ReqCommand) Validate() error {
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

type ResCommand struct {
	Id         string
	Name       string
	Year       int
	Author     string
	Summary    string
	Publisher  string
	PageCount  int
	ReadPage   int
	Finished   bool
	Reading    bool
	IsDeleted  bool
	InsertedAt string
	UpdatedAt  string
}

func mapBookReqCmdToEntity(c ReqCommand) Book {
	e := Book{
		Id:        c.Id,
		Name:      c.Name,
		Year:      c.Year,
		Author:    c.Author,
		Summary:   c.Summary,
		Publisher: c.Publisher,
		PageCount: c.PageCount,
		ReadPage:  c.ReadPage,
		Finished:  c.ReadPage == c.PageCount,
		Reading:   c.Reading,
	}
	return e
}

func mapBookEntityToResCmd(e Book) ResCommand {
	c := ResCommand(e)
	return c
}

func mapBookEntitiesToCmds(es []Book) []ResCommand {
	cs := make([]ResCommand, len(es))
	for i, v := range es {
		cs[i] = mapBookEntityToResCmd(v)
	}
	return cs
}

func mapBookQueryCmdToInfra(c QueryCommand) GetBookQuery {
	i := GetBookQuery(c)
	return i
}

type Service struct {
	Fr *FileRepository
}
type FuncSign func(w http.ResponseWriter, r *http.Request, s *Service)

func (s *Service) Http(f FuncSign) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s)
	}
}

func (s *Service) SaveBook(b ReqCommand) (ResCommand, error) {
	if err := b.Validate(); err != nil {
		return ResCommand{}, fmt.Errorf("not valid: %s", err)
	}

	nb := mapBookReqCmdToEntity(b)
	r := nb.Save(s.Fr)
	c := mapBookEntityToResCmd(r)

	return c, nil
}

func (s *Service) GetBooks(q QueryCommand) []ResCommand {
	var b Book
	bs := b.GetSelectedBooks(s.Fr, mapBookQueryCmdToInfra(q))
	n := 0
	for _, v := range bs {
		if !v.IsDeleted {
			bs[n] = v
			n++
		}
	}

	bs = bs[:n]
	r := mapBookEntitiesToCmds(bs)
	return r
}

func (s *Service) GetBook(id string) (ResCommand, error) {
	var b Book
	bs := b.GetAllBooks(s.Fr)

	for _, v := range bs {
		if v.Id == id && !v.IsDeleted {
			nv := mapBookEntityToResCmd(v)
			return nv, nil
		}
	}

	return ResCommand{}, errors.New("book not Found")
}

func (s *Service) UpdateBook(id string, nb ReqCommand) (ResCommand, error) {
	if err := nb.Validate(); err != nil {
		return ResCommand{}, fmt.Errorf("not valid: %s", err)
	}

	b, err := s.GetBook(id)
	if err != nil {
		return ResCommand{}, err
	}

	nb.Id = b.Id

	eb := mapBookReqCmdToEntity(nb)
	ub, err := eb.Update(s.Fr)
	if err != nil {
		return ResCommand{}, err
	}

	br := mapBookEntityToResCmd(ub)
	return br, nil
}

func (s *Service) DeleteBook(id string) (ResCommand, error) {
	var eb Book
	db, err := eb.Delete(s.Fr, id)
	if err != nil {
		return ResCommand{}, err
	}

	b := mapBookEntityToResCmd(db)
	return b, nil
}
