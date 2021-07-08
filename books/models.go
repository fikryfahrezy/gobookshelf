package books

import (
	"strings"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type BookModel struct {
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

func (b *BookModel) Save() {
	common.Insert(b)
}

func GetAllBooks() []BookModel {
	var b []BookModel
	common.Read(&b)
	return b
}

func GetSelectedBooks(q GetBookQuery) []BookModel {
	b := []BookModel{}
	common.Read(&b)
	var nb []BookModel
	n, f, d := q.Name, q.Finished, q.Reading
	for _, v := range b {
		if d.Exist {
			if d.Val == v.Reading {
				nb = append(nb, v)
			}
		}
		if f.Exist {
			if f.Val == v.Finished {
				nb = append(nb, v)
			}
		}
		if n.Exist {
			if strings.Contains(strings.ToLower(v.Name), strings.ToLower(n.Val)) {
				nb = append(nb, v)
			}
		}
	}
	if nb == nil {
		return b
	}
	return nb
}

func (b *BookModel) Update() {
	bs := GetAllBooks()
	for i, v := range bs {
		if v.Id == b.Id {
			bs[i] = *b
		}
	}
	common.Update(bs)
}

func (b *BookModel) Delete() {
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
