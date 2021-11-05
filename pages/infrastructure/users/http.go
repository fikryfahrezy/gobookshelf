package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
	http_interface "github.com/fikryfahrezy/gobookshelf/pages/interfaces/http"
)

type httpClient struct {
	address string
}

func NewHTTPClient(address string) httpClient {
	return httpClient{address: address}
}

func (h httpClient) Registration(u pages.User) (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/userreg", h.address), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r http_interface.RegResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}
