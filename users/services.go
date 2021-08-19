package users

import (
	"fmt"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/utils"
)

func createUser(nu userModel) (userModel, bool) {
	cu, ok := nu.Save()

	if !ok {
		return userModel{}, false
	}

	return cu, true
}

func getUser(e, p string) (userModel, bool) {
	us, ok := users.ReadByEmail(e)

	if !ok || us.Password != p {
		return userModel{}, false
	}

	return us, true
}

func GetUserById(k string) (userModel, bool) {
	us, ok := users.ReadById(k)

	if !ok {
		return userModel{}, false
	}

	return us, true
}

func updateUser(k string, u userModel) (userModel, bool) {
	c, ok := GetUserById(k)

	if !ok {
		return userModel{}, false
	}

	c, ok = c.Update(u)

	return c, ok
}

func createForgotPass(e string) (string, bool) {
	_, ok := users.ReadByEmail(e)

	if !ok {
		return "", false
	}

	from := "email@email.com"
	code := utils.RandString(15)
	fpM := forgotPassModel{Email: e, Code: code}
	msg := fmt.Sprintf(`
		Code: %s
		<a href="%s/resetpass?code=%s">Click Here</a>
	`, code, common.OwnServerUrl, code)
	err := sendEmail([]string{e}, from, msg)
	if err != nil {
		return err.Error(), false
	}

	fpM.Save()

	return "", true
}

func updateForgotPass(cd, p string) (forgotPassModel, bool) {
	c, ok := ForgotPasses.ReadByCode(cd)

	if !ok {
		return forgotPassModel{}, false
	}

	nfpM, ok := c.Update(forgotPassModel{IsClaimed: true})

	if !ok {
		return forgotPassModel{}, false
	}

	u, ok := users.ReadByEmail(nfpM.Email)

	if !ok {
		return forgotPassModel{}, false
	}

	u, ok = u.Update(userModel{Password: p})

	return nfpM, ok
}
