package user

import (
	"errors"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type User struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func (u *User) Insert(r *Repository) (User, error) {
	u.Id = common.RandString(7)
	for _, v := range r.GetAll() {
		if v.Email == u.Email {
			return User{}, errors.New("email already exist")
		}
	}

	r.Insert(*u)

	return *u, nil
}

func (u *User) ReadByEmail(r *Repository, k string) (User, error) {
	for _, v := range r.GetAll() {
		if v.Email == k {
			return v, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (u *User) ReadById(r *Repository, k string) (User, error) {
	for _, v := range r.GetAll() {
		if v.Id == k {
			return v, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (u *User) Update(r *Repository) (User, error) {

	for _, v := range r.GetAll() {
		if v.Email == u.Email && v.Id != u.Id {
			return User{}, errors.New("user not found")
		}
	}

	for _, v := range r.GetAll() {
		if v.Id == u.Id {
			if u.Email != "" {
				v.Email = u.Email
			}
			if u.Password != "" {
				v.Password = u.Password
			}
			if u.Name != "" {
				v.Name = u.Name
			}
			if u.Region != "" {
				v.Region = u.Region
			}
			if u.Street != "" {
				v.Street = u.Street
			}

			r.UpdateById(v.Id, v)
			return v, nil
		}
	}

	return User{}, errors.New("user not found %s")
}
