package users

import (
	"errors"
	"fmt"

	user_common "github.com/fikryfahrezy/gobookshelf/users"
	"github.com/fikryfahrezy/gobookshelf/users/domain/users"

	"github.com/fikryfahrezy/gobookshelf/users/infrastructure/forgotpw"

	user_repository "github.com/fikryfahrezy/gobookshelf/users/infrastructure/users"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type UserService struct {
	Ur *user_repository.UserRepository
	Fr forgotpw.ForgotPassRepository
}

func (s UserService) CreateUser(nu users.UserModel) (users.UserModel, error) {
	cu, err := s.Ur.Insert(nu)
	if err != nil {
		return users.UserModel{}, err
	}

	return cu, nil
}

func (s UserService) GetUser(e, p string) (users.UserModel, error) {
	us, err := s.Ur.ReadByEmail(e)
	if err != nil {
		return users.UserModel{}, err
	}

	if us.Password != p {
		return users.UserModel{}, errors.New("wrong credential")
	}

	return us, nil
}

func (s UserService) GetUserById(k string) (users.UserModel, error) {
	us, err := s.Ur.ReadById(k)
	if err != nil {
		return users.UserModel{}, err
	}

	return us, nil
}

func (s *UserService) UpdateUser(k string, u users.UserModel) (users.UserModel, error) {
	c, err := s.GetUserById(k)
	if err != nil {
		return users.UserModel{}, err
	}

	c, err = s.Ur.Update(u)

	return c, err
}

func (s UserService) CreateForgotPass(e string) (string, error) {
	_, err := s.Ur.ReadByEmail(e)
	if err != nil {
		return "", err
	}

	code := common.RandString(15)
	fpM := users.ForgotPassModel{Email: e, Code: code}
	from := "email@email.com"
	// msg := fmt.Sprintf(`
	// 	Code: %s
	// 	<a href="%s/resetpass?code=%s">Click Here</a>
	// `, code, handler.OwnServerUrl, code)
	msg := fmt.Sprintf(`
		Code: %s
		<a href="%s/resetpass?code=%s">Click Here</a>
	`, code, "hi", code)

	err = user_common.SendEmail([]string{e}, from, msg)

	if err != nil {
		return "", err
	}

	s.Fr.Insert(fpM)

	return "", nil
}

func (s UserService) GetForgotPass(c string) (users.ForgotPassModel, error) {
	f, err := s.Fr.ReadByCode(c)
	if err != nil {
		return users.ForgotPassModel{}, err
	}

	return f, nil
}

func (s UserService) UpdateForgotPass(cd, p string) (users.ForgotPassModel, error) {
	c, err := s.Fr.ReadByCode(cd)
	if err != nil {
		return users.ForgotPassModel{}, err
	}
	c.IsClaimed = true

	nfpM, err := s.Fr.Update(c)
	if err != nil {
		return users.ForgotPassModel{}, err
	}

	u, err := s.Ur.ReadByEmail(nfpM.Email)
	if err != nil {
		return users.ForgotPassModel{}, err
	}

	u.Password = p
	u, err = s.Ur.Update(u)

	return nfpM, err
}
