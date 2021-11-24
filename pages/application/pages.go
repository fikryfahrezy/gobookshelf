package application

import (
	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
)

type UserReqCommand struct {
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func mapUserCmdToEntity(u UserReqCommand) pages.User {
	us := pages.User{
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

func mapUserEntityToResCmd(u pages.User) UserResCommand {
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

func mapLoginCmdToEntity(l LoginCommand) pages.Auth {
	a := pages.Auth(l)
	return a
}

type ForgotPassResCommand struct {
	Id        int
	Email     string
	Code      string
	IsClaimed bool
}

func mapForgotPassEntityToResCmd(f pages.ForgotPass) ForgotPassResCommand {
	c := ForgotPassResCommand{
		Id:        f.Id,
		Email:     f.Email,
		Code:      f.Code,
		IsClaimed: f.IsClaimed,
	}

	return c
}

type userService interface {
	Registration(u pages.User) (string, error)
	Login(u pages.Auth) (string, error)
	UpdateAcc(a string, u pages.User) (string, error)
	GetUser(a string) (pages.User, error)
	GetForgotPassword(c string) (pages.ForgotPass, error)
}

type galleryService interface {
	GetImages() (interface{}, error)
}

type bookService interface {
	GetBooks(q string) (interface{}, error)
}

type PagesService struct {
	UserService    userService
	GalleryService galleryService
	BookService    bookService
}

func (p PagesService) Registration(u UserReqCommand) (string, error) {
	d, err := p.UserService.Registration(mapUserCmdToEntity(u))
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p PagesService) Login(u LoginCommand) (string, error) {
	d, err := p.UserService.Login(mapLoginCmdToEntity(u))
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p PagesService) UpdateAcc(a string, u UserReqCommand) (string, error) {
	d, err := p.UserService.UpdateAcc(a, mapUserCmdToEntity(u))
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p PagesService) GetUser(a string) (UserResCommand, error) {
	u, err := p.UserService.GetUser(a)
	if err != nil {
		return UserResCommand{}, err
	}

	c := mapUserEntityToResCmd(u)
	return c, nil
}

func (p PagesService) GetForgotPassword(c string) (ForgotPassResCommand, error) {
	s, err := p.UserService.GetForgotPassword(c)
	if err != nil {
		return ForgotPassResCommand{}, err
	}

	f := mapForgotPassEntityToResCmd(s)
	return f, nil

}

func (p PagesService) GetGalleryImages() (interface{}, error) {
	return p.GalleryService.GetImages()
}

func (p PagesService) GetBooks(q string) (interface{}, error) {
	return p.BookService.GetBooks(q)
}
