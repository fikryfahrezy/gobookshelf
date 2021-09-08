package books

import (
	"errors"
)

func saveBook(b bookModel) bookModel {
	b.Save()

	return b
}

func GetBooks(q GetBookQuery) []bookModel {
	bs := GetSelectedBooks(q)
	n := 0

	for _, v := range bs {
		if !v.IsDeleted {
			bs[n] = v
			n++
		}
	}

	bs = bs[:n]

	return bs
}

func getBook(id string) (bookModel, error) {
	b := GetAllBooks()

	for _, v := range b {
		if v.Id == id && !v.IsDeleted {
			return v, nil
		}
	}

	return bookModel{}, errors.New("book not Found")
}

func updateBook(id string, nb bookModel) (bookModel, error) {
	b, err := getBook(id)
	if err != nil {
		return b, err
	}

	b.Update(nb)

	return b, nil
}

func deleteBook(id string) (bookModel, error) {
	b, err := getBook(id)
	if err != nil {
		return b, err
	}

	b.Delete()

	return b, nil
}
