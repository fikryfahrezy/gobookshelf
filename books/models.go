package books

import (
	"strconv"
	"strings"
	"time"

	"github.com/fikryfahrezy/gobookshelf/data"
	"github.com/fikryfahrezy/gobookshelf/utils"
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
	b.Id = utils.RandString(5)
	b.InsertedAt = t
	b.UpdatedAt = t

	data.Insert(b)
}

func (b *bookModel) Update(nb bookModel) {
	var ob bookModel
	bs := GetAllBooks()
	ci := -1

	for i, v := range bs {
		if v.Id == b.Id {
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

		data.Update(bs)
	}
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
