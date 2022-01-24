package page

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthResponseView struct {
	Data string `json:"data"`
}

type UserResponseView struct {
	Data struct {
		Id       string `json:"id"`
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

type UserHttpClient struct {
	Address string
}

func (h UserHttpClient) Registration(u UserReqCommand) (string, error) {
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

func (h UserHttpClient) Login(a Auth) (string, error) {
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

func (h UserHttpClient) UpdateAcc(a string, u UserReqCommand) (string, error) {
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

func (h UserHttpClient) GetUser(a string) (User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/userprofile", h.Address), nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Add("authorization", a)

	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}

	defer resp.Body.Close()

	var r UserResponseView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return User{}, err
	}

	return User(r.Data), nil
}

func (h UserHttpClient) GetForgotPassword(c string) (ForgotPass, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/forgotpassword/%s", h.Address, c), nil)
	if err != nil {
		return ForgotPass{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ForgotPass{}, err
	}

	defer resp.Body.Close()

	var r ForgotPassView
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return ForgotPass{}, err
	}

	fp := ForgotPass{
		Id:        r.Data.Id,
		Email:     r.Data.Email,
		Code:      r.Data.Code,
		IsClaimed: r.Data.IsClaimed,
	}

	return fp, nil
}
