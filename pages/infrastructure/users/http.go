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

	var r http_interface.AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}

func (h httpClient) Login(a pages.Auth) (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/userlogin", h.address), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r http_interface.AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}

func (h httpClient) UpdateAcc(a string, u pages.User) (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/updateuser", h.address), bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", a)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r http_interface.AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}
