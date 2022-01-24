package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fikryfahrezy/gobookshelf/common"
	"net/http"
	"net/mail"
)

type ReqCommand struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func (ur ReqCommand) RegValidate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	if ur.Password == "" {
		return errors.New("password required")
	}

	if ur.Name == "" {
		return errors.New("name required")
	}

	if ur.Region == "" {
		return errors.New("region required")
	}

	if ur.Street == "" {
		return errors.New("street required")
	}

	return nil
}

func (ur ReqCommand) UpValidate() error {
	if _, err := mail.ParseAddress(ur.Email); ur.Email != "" && err != nil {
		return errors.New("input valid email")
	}

	return nil
}

type loginCommand struct {
	Email    string
	Password string
}

func (ur loginCommand) Validate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	if ur.Password == "" {
		return errors.New("password required")
	}

	return nil
}

type forgotPassCommand struct {
	Email string
}

func (ur forgotPassCommand) Validate() error {
	if ur.Email == "" {
		return errors.New("email required")
	} else if _, err := mail.ParseAddress(ur.Email); err != nil {
		return errors.New("input valid email")
	}

	return nil
}

type resetPassCommand struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (ur resetPassCommand) Validate() error {
	if ur.Password == "" {
		return errors.New("password required")
	}

	if ur.Code == "" {
		return errors.New("code required")
	}

	return nil
}

type ResCommand struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func mapUserReqCmdToEntity(u ReqCommand) User {
	nu := User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Region:   u.Region,
		Street:   u.Street,
	}
	return nu
}

func mapUserEntityToResCmd(eu User) ResCommand {
	ur := ResCommand{
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

func mapForgotPwEntityToRes(u ForgotPass) ForgotPwResCommand {
	f := ForgotPwResCommand{
		Id:        u.Id,
		Email:     u.Email,
		Code:      u.Code,
		IsClaimed: u.IsClaimed,
	}

	return f
}

type ForgotPassRepository struct {
	Db *sql.DB
}

type Service struct {
	Ur *Repository
	Fr *ForgotPassRepository
}
type FuncSign func(w http.ResponseWriter, r *http.Request, s *Service)

func (s *Service) Http(f FuncSign) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s)
	}
}

func (s *Service) CreateUser(nu ReqCommand) (ResCommand, error) {
	if err := nu.RegValidate(); err != nil {
		return ResCommand{}, fmt.Errorf("validation fail: %s", err)
	}

	u := mapUserReqCmdToEntity(nu)
	cu, err := u.Insert(s.Ur)
	if err != nil {
		return ResCommand{}, err
	}

	rc := mapUserEntityToResCmd(cu)
	return rc, nil
}

func (s *Service) GetUser(c loginCommand) (ResCommand, error) {
	if err := c.Validate(); err != nil {
		return ResCommand{}, fmt.Errorf("validation fail: %s", err)
	}

	var u User
	us, err := u.ReadByEmail(s.Ur, c.Email)
	if err != nil {
		return ResCommand{}, err
	}

	if us.Password != c.Password {
		return ResCommand{}, errors.New("wrong credential")
	}

	rc := mapUserEntityToResCmd(us)
	return rc, nil
}

func (s *Service) GetUserById(k string) (ResCommand, error) {
	var um User
	us, err := um.ReadById(s.Ur, k)
	if err != nil {
		return ResCommand{}, err
	}

	rc := mapUserEntityToResCmd(us)
	return rc, nil
}

func (s *Service) UpdateUser(k string, u ReqCommand) (ResCommand, error) {
	if err := u.UpValidate(); err != nil {
		return ResCommand{}, fmt.Errorf("validation fail: %s", err)
	}

	var um User
	c, err := um.ReadById(s.Ur, k)
	if err != nil {
		return ResCommand{}, err
	}

	u.Id = c.Id
	ue := mapUserReqCmdToEntity(u)
	c, err = ue.Update(s.Ur)
	if err != nil {
		return ResCommand{}, err
	}

	rc := mapUserEntityToResCmd(c)
	return rc, nil
}

func (s *Service) CreateForgotPass(c forgotPassCommand) (string, error) {
	if err := c.Validate(); err != nil {
		return "", fmt.Errorf("validation fail: %s", err)
	}

	var u User
	_, err := u.ReadByEmail(s.Ur, c.Email)
	if err != nil {
		return "", err
	}

	code := common.RandString(15)
	fpM := ForgotPass{Email: c.Email, Code: code}
	from := "email@email.com"
	// msg := fmt.Sprintf(`
	// 	Code: %s
	// 	<a href="%s/resetpass?code=%s">Click Here</a>
	// `, code, handler.OwnServerUrl, code)
	msg := fmt.Sprintf(`
		Code: %s
		<a href="%s/resetpass?code=%s">Click Here</a>
	`, code, "hi", code)

	err = SendEmail([]string{c.Email}, from, msg)

	if err != nil {
		return "", err
	}

	fpM.Insert(s.Fr)

	return "", nil
}

func (s *Service) GetForgotPass(c string) (ForgotPwResCommand, error) {
	var fr ForgotPass
	u, err := fr.ReadByCode(s.Fr, c)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	f := mapForgotPwEntityToRes(u)
	return f, nil
}

func (s *Service) UpdateForgotPass(cm resetPassCommand) (ForgotPwResCommand, error) {
	if err := cm.Validate(); err != nil {
		return ForgotPwResCommand{}, fmt.Errorf("validation fail: %s", err)
	}

	var fr ForgotPass
	c, err := fr.ReadByCode(s.Fr, cm.Code)
	if err != nil {
		return ForgotPwResCommand{}, err
	}
	c.IsClaimed = true

	nfpM, err := c.Update(s.Fr)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	var um User
	u, err := um.ReadByEmail(s.Ur, nfpM.Email)
	if err != nil {
		return ForgotPwResCommand{}, err
	}

	u.Password = cm.Password
	u, err = u.Update(s.Ur)
	f := mapForgotPwEntityToRes(nfpM)
	return f, err
}
