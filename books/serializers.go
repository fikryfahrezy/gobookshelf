package books

type CommonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type BookIdResponse struct {
	BookId string `json:"bookId"`
}

type book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type BooksSerializer struct {
	Books []BookModel
}

type BooksResponse struct {
	Books []book `json:"books"`
}

type BookSerializer struct {
	Book BookModel
}

func (c *CommonResponse) Response() *CommonResponse {
	return c
}

func (b *BookIdResponse) Response() *BookIdResponse {
	return b
}

func (s *BooksSerializer) Response() BooksResponse {
	var b BooksResponse
	for _, v := range s.Books {
		b.Books = append(b.Books, book{v.Id, v.Name, v.Publisher})
	}
	return b
}
