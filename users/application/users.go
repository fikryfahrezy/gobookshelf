package application

import (
	"errors"
	"fmt"

	user_common "github.com/fikryfahrezy/gobookshelf/users"
	"github.com/fikryfahrezy/gobookshelf/users/domain/users"

	forgotpw_repository "github.com/fikryfahrezy/gobookshelf/users/infrastructure/forgotpw"
	user_repository "github.com/fikryfahrezy/gobookshelf/users/infrastructure/users"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type UserReqCommand struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

type UserResCommand struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func mapUserReqCmdToEntity(u UserReqCommand) users.User {
	nu := users.User{
		Id:       u.Id,
		Email:    u.Name,
		Password: u.Password,
		Name:     u.Name,
		Region:   u.Region,
		Street:   u.Street,
	}
	return nu
}

func mapUserEntityToResCmd(eu users.User) UserResCommand {
	ur := UserResCommand{
		Id:       eu.Id,
		Email:    eu.Email,
		Password: eu.Password,
		Name:     eu.Name,
		Region:   eu.Region,
		Street:   eu.Street,
	}
	return ur
}

type ForgotPwResCommand struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}

func mapForgotPwEntityToRes(u users.ForgotPass) ForgotPwResCommand {
	f := ForgotPwResCommand{
		Id:        u.Id,
		Email:     u.Email,
		Code:      u.Code,
		IsClaimed: u.IsClaimed,
	}

	return f
}

type UserService struct {
	Ur *user_repository.UserRepository
	Fr forgotpw_repository.ForgotPassRepository
}

func (s UserService) CreateUser(nu UserReqCommand) (UserResCommand, error) {
	cu, err := s.Ur.Insert(mapUserReqCmdToEntity(nu))
	if err != nil {
		return UserResCommand{}, err
	}

	rc := mapUserEntityToResCmd(cu)
	return rc, nil
}

func (s UserService) GetUser(e, p string) (UserResCommand, error) {
	us, err := s.Ur.ReadByEmail(e)
	if err != nil {
		return UserResCommand{}, err
	}

	if us.Password != p {
		return UserResCommand{}, errors.New("wrong credential")
	}

	rc := mapUserEntityToResCmd(us)
	return rc, nil
}

func (s UserService) GetUserById(k string) (UserResCommand, error) {
	us, err := s.Ur.ReadById(k)
	if err != nil {
		return UserResCommand{}, err
	}

	rc := mapUserEntityToResCmd(us)
	return rc, nil
}

func (s *UserService) UpdateUser(k string, u UserReqCommand) (UserResCommand, error) {
	c, err := s.Ur.ReadById(k)
	if err != nil {
		return UserResCommand{}, err
	}

	u.Id = c.Id
	c, err = s.Ur.Update(mapUserReqCmdToEntity(u))
	if err != nil {
		return UserResCommand{}, err
	}

	rc := mapUserEntityToResCmd(c)
	return rc, nil
}

func (s UserService) CreateForgotPass(e string) (string, error) {
	_, err := s.Ur.ReadByEmail(e)
	if err != nil {
		return "", err
	}

	code := common.RandString(15)
	fpM := users.ForgotPass{Email: e, Code: code}
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

func (s UserService) GetForgotPass(c string) (ForgotPwResCommand, error) {
	u, err := s.Fr.ReadByCode(c)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	f := mapForgotPwEntityToRes(u)
	return f, nil
}

func (s UserService) UpdateForgotPass(cd, p string) (ForgotPwResCommand, error) {
	c, err := s.Fr.ReadByCode(cd)
	if err != nil {
		return ForgotPwResCommand{}, err
	}
	c.IsClaimed = true

	nfpM, err := s.Fr.Update(c)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	u, err := s.Ur.ReadByEmail(nfpM.Email)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	u.Password = p
	u, err = s.Ur.Update(u)
	f := mapForgotPwEntityToRes(nfpM)
	return f, err
}
