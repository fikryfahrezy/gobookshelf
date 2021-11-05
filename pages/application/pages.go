package application

import (
	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
)

type userService interface {
	Registration(u pages.User) (string, error)
}

type RegistrationCommand struct {
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

type PagesService struct {
	userService userService
}

func (p PagesService) Registration(cmd RegistrationCommand) (string, error) {
	u := pages.User{Email: cmd.Email, Password: cmd.Password, Name: cmd.Name, Region: cmd.Region, Street: cmd.Street}
	d, err := p.userService.Registration(u)
	if err != nil {
		return "", err
	}

	return d, nil
}

func NewPagesServices(uSr userService) PagesService {
	return PagesService{userService: uSr}
}
