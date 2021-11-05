package users

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/handler"
)

func Registration(w http.ResponseWriter, r *http.Request) {
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

	nu := userModel{}
	mapUser(&nu, u)

	ur, err := createUser(nu)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	ur, err := getUser(u.Email, u.Password)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusOK, res)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	nu := userModel{}
	mapUser(&nu, u)

	ur, err := updateUser(c, nu)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ur.Id}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
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

	_, err = createForgotPass(u.Email)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: "Hi"}

	handler.ResJSON(w, http.StatusOK, res.Response())
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	nfpM, err := updateForgotPass(u.Code, u.Pasword)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: nfpM.Id}

	handler.ResJSON(w, http.StatusOK, res.Response())
}
