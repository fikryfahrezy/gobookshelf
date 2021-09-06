package users

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	var u userReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: ""}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.RegValidate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nu := userModel{}
	mapUser(&nu, u)

	ur, err := createUser(nu)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusCreated, res.Response())
}

func Login(w http.ResponseWriter, r *http.Request) {
	var u loginReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: ""}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	ur, err := getUser(u.Email, u.Password)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusOK, res)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var u userReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: ""}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.UpValidate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

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

	ur, err := updateUser(c, nu)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var u forgotPassReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: ""}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	err = createForgotPass(u.Email)

	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: "Hi"}

	common.ResJSON(w, http.StatusOK, res.Response())
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u resetPassReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: ""}

		common.ResJSON(w, errDcd.Status, res.Response())
		return
	}

	err := u.Validate()
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nfpM, err := updateForgotPass(u.Code, u.Pasword)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: nfpM.Id}

	common.ResJSON(w, http.StatusOK, res.Response())
}
