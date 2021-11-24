package users

import (
	"errors"
	"sync"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/users/domain/users"
)

type UserRepository struct {
	Users map[string]users.User
	lock  sync.RWMutex
}

func (udb *UserRepository) Insert(u users.User) (users.User, error) {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	u.Id = common.RandString(7)
	for _, v := range udb.Users {
		if v.Email == u.Email {
			return users.User{}, errors.New("email already exist")
		}
	}

	udb.Users[common.RandString(9)] = u

	return u, nil
}

func (udb *UserRepository) ReadByEmail(k string) (users.User, error) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.Users {
		if v.Email == k {
			return v, nil
		}
	}

	return users.User{}, errors.New("user not found")
}

func (udb *UserRepository) ReadById(k string) (users.User, error) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.Users {
		if v.Id == k {
			return v, nil
		}
	}

	return users.User{}, errors.New("user not found")
}

func (udb *UserRepository) Update(u users.User) (users.User, error) {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	for _, v := range udb.Users {
		if v.Email == u.Email && v.Id != u.Id {
			return users.User{}, errors.New("user not foundx")
		}
	}

	for i, v := range udb.Users {
		if v.Id == u.Id {
			ou := udb.Users[i]

			if u.Email != "" {
				ou.Email = u.Email
			}
			if u.Password != "" {
				ou.Password = u.Password
			}
			if u.Name != "" {
				ou.Name = u.Name
			}
			if u.Region != "" {
				ou.Region = u.Region
			}
			if u.Street != "" {
				ou.Street = u.Street
			}

			udb.Users[i] = u

			return udb.Users[i], nil
		}
	}

	return users.User{}, errors.New("user not found %s")
}
