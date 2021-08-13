package books

import (
	"strconv"
	"strings"
	"time"

	"github.com/fikryfahrezy/gobookshelf/data"
)

type bookModel struct {
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
	InsertedAt string `json:"insertedAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func (b *bookModel) Save() {
	t := time.Now().UTC().String()
	b.Id = t
	b.InsertedAt = t
	b.UpdatedAt = t

	data.Insert(b)
}

func GetAllBooks() []bookModel {
	var b []bookModel
	data.Read(&b)

	return b
}

func GetSelectedBooks(q GetBookQuery) []bookModel {
	b := []bookModel{}
	data.Read(&b)

	var nb []bookModel
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

func (b *bookModel) Update() {
	b.UpdatedAt = time.Now().UTC().String()
	bs := GetAllBooks()

	for i, v := range bs {
		if v.Id == b.Id {
			bs[i] = *b
		}
	}

	data.Update(bs)
}

func (b *bookModel) Delete() {
	bs := GetAllBooks()

	for i, v := range bs {
		if v.Id == b.Id {
			l := len(bs) - 1
			bs[i] = bs[l]
			bs = bs[:l]
		}
	}

	data.Update(bs)
}
