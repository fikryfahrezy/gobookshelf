package page

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BookHttpClient struct {
	Address string
}

type book struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Publisher string `json:"publisher"`
}

type booksRes struct {
	Books []book `json:"books"`
}

type BookClientRes struct {
	Data booksRes `json:"data"`
}

func (h BookHttpClient) GetBooks(q string) (BookClientRes, error) {

	resp, err := http.Get(fmt.Sprintf("%s/books?name=%s", h.Address, q))
	if err != nil {
		return BookClientRes{}, err
	}

	defer resp.Body.Close()

	var r BookClientRes
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return BookClientRes{}, err
	}

	return r, nil
}
