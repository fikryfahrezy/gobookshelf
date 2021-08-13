package users

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	res := commonResponse{}
	var u regReqValidator
	err := common.DecodeJSONBody(w, r, &u)
	if err != nil {

		common.ResJSON(w, err.Status, res.Response())
		return
	}

	_, ok := u.Validate()

	if !ok {

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	nu := userModel{}

	mapUser(&nu, u)
	createUser(nu)

	common.ResJSON(w, http.StatusCreated, res)
}
