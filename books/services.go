package books

import (
	"time"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func saveBook(b *BookModelValidator) BookModel {
	t := time.Now().UTC().String()
	nb := BookModel{
		common.RandString(7),
		b.Name,
		b.Year,
		b.Author,
		b.Summary,
		b.Publisher,
		b.PageCount,
		b.ReadPage,
		b.PageCount == b.ReadPage,
		b.Reading,
		t,
		t,
	}
	nb.Save()

	return nb
}

func getBooks() []BookModel {
	b := GetBooks()
	return b
}
