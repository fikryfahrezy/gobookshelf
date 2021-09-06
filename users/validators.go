package users

import (
	"errors"
	"net/mail"
)

type userReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func (ur *userReq) RegValidate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	if ur.Password == "" {
		return errors.New("password required")
	}

	if ur.Name == "" {
		return errors.New("name required")
	}

	if ur.Region == "" {
		return errors.New("region required")
	}

	if ur.Street == "" {
		return errors.New("street required")
	}

	return nil
}

func (ur *userReq) UpValidate() error {
	if _, err := mail.ParseAddress(ur.Email); ur.Email != "" && err != nil {
		return errors.New("input valid email")
	}

	return nil
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ur *loginReq) Validate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	if ur.Password == "" {
		return errors.New("password required")
	}

	return nil
}

type forgotPassReq struct {
	Email string `json:"email"`
}

func (ur *forgotPassReq) Validate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	return nil
}

type resetPassReq struct {
	Code    string `json:"code"`
	Pasword string `json:"password"`
}

func (ur *resetPassReq) Validate() error {
	if ur.Pasword == "" {
		return errors.New("password required")
	}

	if ur.Code == "" {
		return errors.New("code required")
	}

	return nil
}
