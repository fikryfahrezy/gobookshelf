package users

import (
	"errors"
	"fmt"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/utils"
)

func createUser(nu userModel) (userModel, error) {
	cu, err := nu.Save()
	if err != nil {
		return userModel{}, err
	}

	return cu, nil
}

func getUser(e, p string) (userModel, error) {
	us, err := users.ReadByEmail(e)
	if err != nil {
		return userModel{}, err
	}

	if us.Password != p {
		return userModel{}, errors.New("wrong credential")
	}

	return us, nil
}

func GetUserById(k string) (userModel, error) {
	us, err := users.ReadById(k)
	if err != nil {
		return userModel{}, err
	}

	return us, nil
}

func updateUser(k string, u userModel) (userModel, error) {
	c, err := GetUserById(k)
	if err != nil {
		return userModel{}, err
	}

	c, err = c.Update(u)

	return c, err
}

func createForgotPass(e string) error {
	_, err := users.ReadByEmail(e)
	if err != nil {
		return err
	}

	from := "email@email.com"
	code := utils.RandString(15)
	fpM := forgotPassModel{Email: e, Code: code}
	msg := fmt.Sprintf(`
		Code: %s
		<a href="%s/resetpass?code=%s">Click Here</a>
	`, code, common.OwnServerUrl, code)
	err = sendEmail([]string{e}, from, msg)

	if err != nil {
		return err
	}

	fpM.Save()

	return nil
}

func updateForgotPass(cd, p string) (forgotPassModel, error) {
	c, err := ForgotPasses.ReadByCode(cd)
	if err != nil {
		return forgotPassModel{}, err
	}

	nfpM, err := c.Update(forgotPassModel{IsClaimed: true})
	if err != nil {
		return forgotPassModel{}, err
	}

	u, err := users.ReadByEmail(nfpM.Email)
	if err != nil {
		return forgotPassModel{}, err
	}

	u, err = u.Update(userModel{Password: p})

	return nfpM, err
}
