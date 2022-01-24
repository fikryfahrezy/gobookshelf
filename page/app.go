package page

import (
	"html/template"
	"net/http"
)

type UserReqCommand struct {
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func mapUserCmdToEntity(u UserReqCommand) User {
	us := User{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Region:   u.Region,
		Street:   u.Street,
	}
	return us
}

type UserResCommand struct {
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func mapUserEntityToResCmd(u User) UserResCommand {
	us := UserResCommand{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Region:   u.Region,
		Street:   u.Street,
	}
	return us
}

type LoginCommand struct {
	Email    string
	Password string
}

func mapLoginCmdToEntity(l LoginCommand) Auth {
	a := Auth(l)
	return a
}

type ForgotPassResCommand struct {
	Id        int
	Email     string
	Code      string
	IsClaimed bool
}

func mapForgotPassEntityToResCmd(f ForgotPass) ForgotPassResCommand {
	c := ForgotPassResCommand{
		Id:        f.Id,
		Email:     f.Email,
		Code:      f.Code,
		IsClaimed: f.IsClaimed,
	}

	return c
}

type UserService interface {
	Registration(u UserReqCommand) (string, error)
	Login(u Auth) (string, error)
	UpdateAcc(a string, u UserReqCommand) (string, error)
	GetUser(a string) (User, error)
	GetForgotPassword(c string) (ForgotPass, error)
}

type GalleryService interface {
	GetImages() (ImageClientRes, error)
}

type BookService interface {
	GetBooks(q string) (BookClientRes, error)
}

type Service struct {
	Session  *UserSession
	Template *template.Template
	Host     string
	UserService
	GalleryService
	BookService
}

type FuncSign func(w http.ResponseWriter, r *http.Request, s *Service)

func (p *Service) Http(f FuncSign) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, p)
	}
}

func (p *Service) UserRegistration(u UserReqCommand) (string, error) {
	d, err := p.Registration(u)
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p *Service) UserLogin(u LoginCommand) (string, error) {
	d, err := p.Login(mapLoginCmdToEntity(u))
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p *Service) UpdateUserAcc(a string, u UserReqCommand) (string, error) {
	d, err := p.UpdateAcc(a, u)
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p *Service) GetUserProfile(a string) (UserResCommand, error) {
	u, err := p.GetUser(a)
	if err != nil {
		return UserResCommand{}, err
	}

	c := mapUserEntityToResCmd(u)
	return c, nil
}

func (p *Service) GetUserForgotPassword(c string) (ForgotPassResCommand, error) {
	s, err := p.GetForgotPassword(c)
	if err != nil {
		return ForgotPassResCommand{}, err
	}

	f := mapForgotPassEntityToResCmd(s)
	return f, nil

}
