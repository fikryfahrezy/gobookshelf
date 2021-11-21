package application

import (
	"errors"

	book_domain "github.com/fikryfahrezy/gobookshelf/books/domain/books"

	"github.com/fikryfahrezy/gobookshelf/books/infrastructure/books"
)

type GetBookQuery struct {
	Name     string
	Reading  string
	Finished string
}

type BookService struct {
	Fr books.FileRepository
}

func (s BookService) SaveBook(b book_domain.BookModel) book_domain.BookModel {
	return s.Fr.Save(b)
}

func (s BookService) GetBooks(q GetBookQuery) []book_domain.BookModel {
	bs := s.Fr.GetSelectedBooks(books.GetBookQuery{
		Name:     q.Name,
		Reading:  q.Reading,
		Finished: q.Finished,
	})

	n := 0
	for _, v := range bs {
		if !v.IsDeleted {
			bs[n] = v
			n++
		}
	}
	bs = bs[:n]

	return bs
}

func (s BookService) GetBook(id string) (book_domain.BookModel, error) {
	b := s.Fr.GetAllBooks()

	for _, v := range b {
		if v.Id == id && !v.IsDeleted {
			return v, nil
		}
	}

	return book_domain.BookModel{}, errors.New("book not Found")
}

func (s BookService) UpdateBook(id string, nb book_domain.BookModel) (book_domain.BookModel, error) {
	b, err := s.GetBook(id)
	if err != nil {
		return book_domain.BookModel{}, err
	}

	nb.Id = b.Id
	b, err = s.Fr.Update(nb)
	if err != nil {
		return book_domain.BookModel{}, err
	}

	return b, nil
}

func (s BookService) DeleteBook(id string) (book_domain.BookModel, error) {
	db, err := s.Fr.Delete(id)
	if err != nil {
		return book_domain.BookModel{}, err
	}

	return db, nil
}
