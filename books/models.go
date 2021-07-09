package books

import (
	"strconv"
	"strings"

	"github.com/fikryfahrezy/gobookshelf/common"
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
	common.Insert(b)
}

func GetAllBooks() []bookModel {
	var b []bookModel
	common.Read(&b)
	return b
}

func GetSelectedBooks(q getBookQuery) []bookModel {
	b := []bookModel{}
	common.Read(&b)
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
	bs := GetAllBooks()
	for i, v := range bs {
		if v.Id == b.Id {
			bs[i] = *b
		}
	}
	common.Update(bs)
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
	common.Update(bs)
}
