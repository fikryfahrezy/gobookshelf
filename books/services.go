package books

import (
	"time"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

func newBook(o string, ob *bookModel, nb bookModelValidator) {
	t := time.Now().UTC().String()
	ob.Name = nb.Name
	ob.Year = nb.Year
	ob.Author = nb.Author
	ob.Summary = nb.Summary
	ob.Publisher = nb.Publisher
	ob.PageCount = nb.PageCount
	ob.ReadPage = nb.ReadPage
	ob.Reading = nb.Reading
	ob.Finished = nb.ReadPage == nb.PageCount
	ob.UpdatedAt = t
	if o == "CREATE" {
		ob.Id = utils.RandString(5)
		ob.InsertedAt = t
	}
}

func saveBook(b bookModelValidator) bookModel {
	nb := bookModel{}
	newBook("CREATE", &nb, b)
	nb.Save()
	return nb
}

func GetBooks(q GetBookQuery) []bookModel {
	b := GetSelectedBooks(q)
	return b
}

func getBook(id string) (bookModel, bool) {
	b := GetAllBooks()
	for _, v := range b {
		if v.Id == id {
			return v, true
		}
	}
	return bookModel{}, false
}

func updateBook(id string, nb bookModelValidator) (bookModel, bool) {
	b, ok := getBook(id)
	if !ok {
		return b, ok
	}
	newBook("UPDATE", &b, nb)
	b.Update()
	return b, true
}

func deleteBook(id string) (bookModel, bool) {
	b, ok := getBook(id)
	if !ok {
		return b, ok
	}
	b.Delete()
	return b, true
}
