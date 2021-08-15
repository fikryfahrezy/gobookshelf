package users

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/utils"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	var u regReqValidator
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	msg, ok := u.Validate()

	if !ok {
		res := common.CommonResponse{Status: "fail", Message: msg, Data: make([]interface{}, 0)}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nu := userModel{}

	mapUser(&nu, u)
	createUser(nu)

	res := common.CommonResponse{Status: "success", Message: "", Data: make([]interface{}, 0)}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: utils.RandString(15), HttpOnly: true, Secure: true, SameSite: 3})
	common.ResJSON(w, http.StatusCreated, res)
}
