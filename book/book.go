package book

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type Book struct {
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

type GetBookQuery struct {
	Name     string
	Reading  string
	Finished string
}

func (b *Book) Save(r *FileRepository) Book {
	t := time.Now().UTC().String()
	b.Id = common.RandString(5)
	b.InsertedAt = t
	b.UpdatedAt = t

	r.Insert(b)

	return *b
}

func (b *Book) Update(r *FileRepository) (Book, error) {
	var bs []Book
	r.Read(&bs)

	var ob Book
	ci := -1

	for i, v := range bs {
		if v.Id == b.Id {
			ci = i
			ob = v
			break
		}
	}

	if ci >= 0 {
		ob.Name = b.Name
		ob.Year = b.Year
		ob.Author = b.Author
		ob.Summary = b.Summary
		ob.Publisher = b.Publisher
		ob.PageCount = b.PageCount
		ob.ReadPage = b.ReadPage
		ob.Reading = b.Reading
		ob.Finished = b.ReadPage == b.PageCount
		ob.UpdatedAt = time.Now().UTC().String()

		bs[ci] = ob
		r.Rewrite(bs)

		return ob, nil
	}

	return Book{}, errors.New("not found")
}

func (b *Book) Delete(r *FileRepository, id string) (Book, error) {
	var bs []Book
	r.Read(&bs)

	for i, v := range bs {
		if v.Id == id {
			bs[i].IsDeleted = true
			r.Rewrite(bs)
			return bs[i], nil
		}
	}

	return Book{}, errors.New("not found")
}

func (b *Book) GetAllBooks(r *FileRepository) []Book {
	var bs []Book
	r.Read(&bs)

	return bs
}

func (b *Book) GetSelectedBooks(r *FileRepository, q GetBookQuery) []Book {
	var bs []Book
	r.Read(&bs)

	var nb []Book
	n, f, d := q.Name, q.Finished, q.Reading

	for _, v := range bs {
		if pf, err := strconv.ParseBool(f); err == nil && pf == v.Finished {
			nb = append(nb, v)
		}
		if pd, err := strconv.ParseBool(d); err == nil && pd == v.Reading {
			nb = append(nb, v)
		}
		if n != "" && strings.Contains(strings.ToLower(v.Name), strings.ToLower(n)) {
			nb = append(nb, v)
		}
	}

	if nb == nil {
		return bs
	}

	return nb
}
