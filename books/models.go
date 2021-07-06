package books

import (
	"github.com/fikryfahrezy/gobookshelf/common"
)

type BookModel struct {
	Name      string
	Year      int
	Author    string
	Summary   string
	Publisher string
	PageCount int
	ReadPage  int
	Reading   bool
}

func (b *BookModel) Save() {
	common.Insert(b)
}
