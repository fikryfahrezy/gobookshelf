package books

import (
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
	UpdatedAt  string `json:"udatedAt"`
}

func (b *BookModel) Save() {
	common.Insert(b)
}

func GetBooks() []BookModel {
	var b []BookModel
	common.Read(&b)
	return b
}
