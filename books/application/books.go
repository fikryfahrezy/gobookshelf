package application

import (
	"errors"

	book_domain "github.com/fikryfahrezy/gobookshelf/books/domain/books"

	"github.com/fikryfahrezy/gobookshelf/books/infrastructure/books"
)

type BookQueryCommand struct {
	Name     string
	Reading  string
	Finished string
}

type BookReqCommand struct {
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

type BookResCommand struct {
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

func mapBookReqCmdToEntity(c BookReqCommand) book_domain.Book {
	e := book_domain.Book{
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

func mapBookEntityToResCmd(e book_domain.Book) BookResCommand {
	c := BookResCommand(e)
	return c
}

func mapBookEntitiesToCmds(es []book_domain.Book) []BookResCommand {
	cs := make([]BookResCommand, len(es))
	for i, v := range es {
		cs[i] = mapBookEntityToResCmd(v)
	}
	return cs
}

func mapBookQueryCmdToInfra(c BookQueryCommand) books.GetBookQuery {
	i := books.GetBookQuery(c)
	return i
}

type BookService struct {
	Fr books.FileRepository
}

func (s BookService) SaveBook(b BookReqCommand) BookResCommand {
	r := s.Fr.Save(mapBookReqCmdToEntity(b))
	c := mapBookEntityToResCmd(r)
	return c
}

func (s BookService) GetBooks(q BookQueryCommand) []BookResCommand {
	bs := s.Fr.GetSelectedBooks(mapBookQueryCmdToInfra(q))
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

func (s BookService) GetBook(id string) (BookResCommand, error) {
	b := s.Fr.GetAllBooks()

	for _, v := range b {
		if v.Id == id && !v.IsDeleted {
			nv := mapBookEntityToResCmd(v)
			return nv, nil
		}
	}

	return BookResCommand{}, errors.New("book not Found")
}

func (s BookService) UpdateBook(id string, nb BookReqCommand) (BookResCommand, error) {
	b, err := s.GetBook(id)
	if err != nil {
		return BookResCommand{}, err
	}

	nb.Id = b.Id
	ub, err := s.Fr.Update(mapBookReqCmdToEntity(nb))
	if err != nil {
		return BookResCommand{}, err
	}

	br := mapBookEntityToResCmd(ub)
	return br, nil
}

func (s BookService) DeleteBook(id string) (BookResCommand, error) {
	db, err := s.Fr.Delete(id)
	if err != nil {
		return BookResCommand{}, err
	}

	b := mapBookEntityToResCmd(db)
	return b, nil
}
