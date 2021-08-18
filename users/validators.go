package users

import "net/mail"

type userReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func (ur *userReq) RegValidate() (string, bool) {
	if ur.Email == "" {
		return "email required", false
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return "input valid email", false
	}

	if ur.Password == "" {
		return "password required", false
	}

	if ur.Name == "" {
		return "name required", false
	}

	if ur.Region == "" {
		return "region required", false
	}

	if ur.Street == "" {
		return "street required", false
	}

	return "", true
}

func (ur *userReq) UpValidate() (string, bool) {
	if _, err := mail.ParseAddress(ur.Email); ur.Email != "" && err != nil {
		return "input valid email", false
	}

	return "", true
}

type loginReqValidator struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ur *loginReqValidator) Validate() (string, bool) {
	if ur.Email == "" {
		return "email required", false
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return "input valid email", false
	}

	if ur.Password == "" {
		return "password required", false
	}

	return "", true
}
