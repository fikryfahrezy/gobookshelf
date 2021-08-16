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

	us := UserSessions.Create(ur.Id)
	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	http.SetCookie(w, &http.Cookie{Name: AuthSessionKey, Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	common.ResJSON(w, http.StatusCreated, res)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var u loginReqValidator
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

	us := UserSessions.Create(ur.Id)
	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}

	http.SetCookie(w, &http.Cookie{Name: AuthSessionKey, Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	common.ResJSON(w, http.StatusOK, res)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(AuthSessionKey)
	if err != nil {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	UserSessions.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{Name: AuthSessionKey, MaxAge: -1})
	http.Redirect(w, r, "/home", http.StatusFound)
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

	c, ec := r.Cookie(AuthSessionKey)

	if ec != nil {
		res := common.CommonResponse{Status: "fail", Message: ec.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	uc := UserSessions.Get(c.Value)
	nu := userModel{}
	mapUser(&nu, u)

	ur, ok := updateUser(uc, nu)

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: "", Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ur.Id}
	common.ResJSON(w, http.StatusOK, res)
}
