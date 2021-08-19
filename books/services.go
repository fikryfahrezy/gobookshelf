package books

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

func getBook(id string) (bookModel, bool) {
	b := GetAllBooks()

	for _, v := range b {
		if v.Id == id && !v.IsDeleted {
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

	b.Update(nb)

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
