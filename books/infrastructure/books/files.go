package books

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fikryfahrezy/gobookshelf/books/domain/books"
	"github.com/fikryfahrezy/gobookshelf/common"
)

type GetBookQuery struct {
	Name     string
	Reading  string
	Finished string
}

type FileRepository struct {
	Filename string
}

func (r FileRepository) WriteFile(b []byte) {
	fo, err := os.Create(r.Filename)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := fo.Write(b); err != nil {
		panic(err)
	}
}

func (r FileRepository) Insert(v interface{}) {
	fi, err := os.Open(r.Filename)
	if err != nil {
		panic(err)
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var d []interface{}
	json.NewDecoder(fi).Decode(&d)

	d = append(d, v)
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	r.WriteFile(b)
}

func (r FileRepository) Read(v interface{}) {
	fi, err := os.Open(r.Filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	json.NewDecoder(fi).Decode(&v)
}

func (r FileRepository) Rewrite(v interface{}) {
	fi, err := os.Open(r.Filename)
	if err != nil {
		panic(err)
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	r.WriteFile(b)
}

func (r FileRepository) Save(b books.BookModel) books.BookModel {
	t := time.Now().UTC().String()
	b.Id = common.RandString(5)
	b.InsertedAt = t
	b.UpdatedAt = t

	r.Insert(b)

	return b
}

func (r FileRepository) Update(nb books.BookModel) (books.BookModel, error) {
	var ob books.BookModel
	bs := r.GetAllBooks()
	ci := -1

	for i, v := range bs {
		if v.Id == nb.Id {
			ci = i
			ob = v
			break
		}
	}

	if ci >= 0 {
		ob.Name = nb.Name
		ob.Year = nb.Year
		ob.Author = nb.Author
		ob.Summary = nb.Summary
		ob.Publisher = nb.Publisher
		ob.PageCount = nb.PageCount
		ob.ReadPage = nb.ReadPage
		ob.Reading = nb.Reading
		ob.Finished = nb.ReadPage == nb.PageCount
		ob.UpdatedAt = time.Now().UTC().String()

		bs[ci] = ob

		r.Rewrite(bs)

		return ob, nil
	}

	return books.BookModel{}, errors.New("not found")
}

func (r FileRepository) Delete(id string) (books.BookModel, error) {
	bs := r.GetAllBooks()

	for i, v := range bs {
		if v.Id == id {
			bs[i].IsDeleted = true
			r.Rewrite(bs)
			return bs[i], nil
		}
	}

	return books.BookModel{}, errors.New("not found")
}

func (r FileRepository) GetAllBooks() []books.BookModel {
	var b []books.BookModel
	r.Read(&b)

	return b
}

func (r FileRepository) GetSelectedBooks(q GetBookQuery) []books.BookModel {
	var b []books.BookModel
	r.Read(&b)

	var nb []books.BookModel
	n, f, d := q.Name, q.Finished, q.Reading

	for _, v := range b {
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
		return b
	}

	return nb
}

// InitDB
//
// How to read/write from/to a file using Go
// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
func InitDB(f string) FileRepository {
	fl := FileRepository{Filename: f}

	// open input file
	if _, err := os.Stat(f); os.IsNotExist(err) {
		fl.WriteFile([]byte("[]"))
	}

	return fl
}
