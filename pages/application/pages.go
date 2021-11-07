package application

import (
	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
)

type userService interface {
	Registration(u pages.User) (string, error)
	Login(u pages.Auth) (string, error)
	UpdateAcc(a string, u pages.User) (string, error)
}

type UserCommand struct {
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

type LoginCommand struct {
	Email    string
	Password string
}

type PagesService struct {
	userService userService
}

func (p PagesService) Registration(u UserCommand) (string, error) {
	us := pages.User(u)
	d, err := p.userService.Registration(us)
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p PagesService) Login(u LoginCommand) (string, error) {
	a := pages.Auth(u)
	d, err := p.userService.Login(a)
	if err != nil {
		return "", err
	}

	return d, nil
}

func (p PagesService) UpdateAcc(a string, u UserCommand) (string, error) {
	us := pages.User(u)
	d, err := p.userService.UpdateAcc(a, us)
	if err != nil {
		return "", err
	}

	return d, nil
}

func NewPagesServices(uSr userService) PagesService {
	return PagesService{userService: uSr}
}
