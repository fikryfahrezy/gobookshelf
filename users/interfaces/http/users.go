package http

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/fikryfahrezy/gosrouter"

	"github.com/fikryfahrezy/gobookshelf/users/domain/users"

	user_service "github.com/fikryfahrezy/gobookshelf/users/application"

	"github.com/fikryfahrezy/gobookshelf/handler"
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

type UserView struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func mapUser(um *users.UserModel, ur userReq) {
	um.Email = ur.Email
	um.Region = ur.Region
	um.Street = ur.Street
	um.Name = ur.Name
	um.Password = ur.Password
}

type UserRoutes struct {
	Us user_service.UserService
}

func (s *UserRoutes) Registration(w http.ResponseWriter, r *http.Request) {
	var u userReq
	errDcd := handler.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: ""}

		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.RegValidate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nu := users.UserModel{}
	mapUser(&nu, u)

	ur, err := s.Us.CreateUser(nu)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func (s *UserRoutes) Login(w http.ResponseWriter, r *http.Request) {
	var u loginReq
	errDcd := handler.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: ""}

		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	ur, err := s.Us.GetUser(u.Email, u.Password)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusOK, res)
}

func (s *UserRoutes) GetProfile(w http.ResponseWriter, r *http.Request) {
	c := r.Header.Get("authorization")

	if c == "" {
		res := handler.CommonResponse{Message: http.StatusText(http.StatusUnauthorized), Data: ""}

		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	ur, err := s.Us.GetUserById(c)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (s *UserRoutes) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var u userReq
	errDcd := handler.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: ""}

		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.UpValidate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	c := r.Header.Get("authorization")

	if c == "" {
		res := handler.CommonResponse{Message: http.StatusText(http.StatusUnauthorized), Data: ""}

		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	nu := users.UserModel{}
	mapUser(&nu, u)

	ur, err := s.Us.UpdateUser(c, nu)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (s *UserRoutes) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var u forgotPassReq
	errDcd := handler.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: ""}

		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	_, err = s.Us.CreateForgotPass(u.Email)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: "Hi"}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (s *UserRoutes) GetForgotPassword(w http.ResponseWriter, r *http.Request) {
	p := gosrouter.ReqParams(r.URL.Path)
	c := p("code")

	f, err := s.Us.GetForgotPass(c)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: f}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func (s *UserRoutes) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u resetPassReq
	errDcd := handler.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Error(), Data: ""}

		handler.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nfpM, err := s.Us.UpdateForgotPass(u.Code, u.Pasword)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: nfpM.Id}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func AddRoutes(u UserRoutes) {
	gosrouter.HandlerPOST("/userreg", u.Registration)
	gosrouter.HandlerPOST("/profile", u.GetProfile)
	gosrouter.HandlerPOST("/userlogin", u.Login)
	gosrouter.HandlerPATCH("/updateuser", u.UpdateProfile)
	gosrouter.HandlerPOST("/forgotpassword", u.ForgotPassword)
	gosrouter.HandlerGET("/forgotpassword/:code", u.GetForgotPassword)
	gosrouter.HandlerPATCH("/updatepassword", u.UpdatePassword)
}
