package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
)

type AuthResponseView struct {
	Data string `json:"data"`
}

type UserResponseView struct {
	Data struct {
		Id       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"-"`
		Name     string `json:"name"`
		Region   string `json:"region"`
		Street   string `json:"street"`
	} `json:"data"`
}

type ForgotPassView struct {
	Data struct {
		Id        int    `json:"id"`
		Email     string `json:"email"`
		Code      string `json:"code"`
		IsClaimed bool   `json:"isclaimed"`
	} `json:"data"`
}

type HttpClient struct {
	Address string
}

func (h HttpClient) Registration(u pages.User) (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/userreg", h.Address), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}

func (h HttpClient) Login(a pages.Auth) (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/userlogin", h.Address), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}

func (h HttpClient) UpdateAcc(a string, u pages.User) (string, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/updateprofile", h.Address), bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", a)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r AuthResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Data, nil
}

func (h HttpClient) GetUser(a string) (pages.User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile", h.Address), nil)
	if err != nil {
		return pages.User{}, err
	}

	req.Header.Add("authorization", a)

	resp, err := client.Do(req)
	if err != nil {
		return pages.User{}, err
	}

	defer resp.Body.Close()

	var r UserResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return pages.User{}, err
	}

	return pages.User(r.Data), nil
}

func (h HttpClient) GetForgotPassword(c string) (pages.ForgotPass, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/forgotpassword/%s", h.Address, c), nil)
	if err != nil {
		return pages.ForgotPass{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return pages.ForgotPass{}, err
	}

	defer resp.Body.Close()

	var r ForgotPassView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return pages.ForgotPass{}, err
	}

	fp := pages.ForgotPass{
		Id:        r.Data.Id,
		Email:     r.Data.Email,
		Code:      r.Data.Code,
		IsClaimed: r.Data.IsClaimed,
	}

	return fp, nil
}
