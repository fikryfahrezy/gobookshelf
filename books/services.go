package books

func saveBook(b bookModel) bookModel {
	b.Save()

	return b
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

func updateBook(id string, nb bookModel) (bookModel, bool) {
	b, ok := getBook(id)

	if !ok {
		return b, ok
	}

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
