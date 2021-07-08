package books

type CommonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *CommonResponse) Response() *CommonResponse {
	return c
}

type BookIdResponse struct {
	BookId string `json:"bookId"`
}

func (b *BookIdResponse) Response() *BookIdResponse {
	return b
}

type Book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type BooksSerializer struct {
	Books []BookModel
}

type BooksResponse struct {
	Books []Book `json:"books"`
}

func (s *BooksSerializer) Response() BooksResponse {
	b := BooksResponse{}
	for _, v := range s.Books {
		b.Books = append(b.Books, Book{v.Id, v.Name, v.Publisher})
	}
	return b
}

type BookSerializer struct {
	Book BookModel `json:"book"`
}

func (b *BookSerializer) Response() BookSerializer {
	return *b
}
