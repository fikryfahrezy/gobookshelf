package users

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	var u userReq
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.RegValidate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nu := userModel{}
	mapUser(&nu, u)

	ur, ok := createUser(nu)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusCreated, res.Response())
}

func Login(w http.ResponseWriter, r *http.Request) {
	var u loginReq
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	ur, ok := getUser(u.Email, u.Password)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusOK, res)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var u userReq
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.UpValidate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	c := r.Header.Get("authorization")

	if c == "" {
		res := common.CommonResponse{Status: "fail", Message: http.StatusText(http.StatusUnauthorized), Data: ""}

		common.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	nu := userModel{}
	mapUser(&nu, u)

	ur, ok := updateUser(c, nu)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "", Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var u forgotPassReq
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	msg, ok = createForgotPass(u.Email)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: "Hi"}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u resetPassReq
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nfpM, ok := updateForgotPass(u.Code, u.Pasword)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: ""}

		common.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: nfpM.Id}

	common.ResJSON(w, http.StatusOK, res.Response())
}
