package books

type bookIdResponse struct {
	BookId string `json:"bookId"`
}

func (b *bookIdResponse) Response() *bookIdResponse {
	return b
}

type book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type BooksSerializer struct {
	Books []bookModel
}

type booksResponse struct {
	Books []book `json:"books"`
}

func (s *BooksSerializer) Response() booksResponse {
	b := booksResponse{}

	for _, v := range s.Books {
		b.Books = append(b.Books, book{v.Id, v.Name, v.Publisher})
	}

	if b.Books == nil {
		b.Books = make([]book, 0)
	}

	return b
}

type BookSerializer struct {
	Book bookModel `json:"book"`
}

func (b *BookSerializer) Response() BookSerializer {
	return *b
}
