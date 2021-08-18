package users

import (
	"sync"
	"time"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

type userModel struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func (um *userModel) Save() (userModel, bool) {
	um.Id = utils.RandString(7)
	ur := users.Insert(*um)

	return *um, ur
}

func (um *userModel) Update(nu userModel) (userModel, bool) {
	if nu.Email != "" {
		um.Email = nu.Email
	}

	if nu.Password != "" {
		um.Password = nu.Password
	}

	if nu.Name != "" {
		um.Name = nu.Name
	}

	if nu.Region != "" {
		um.Region = nu.Region
	}

	if nu.Street != "" {
		um.Street = nu.Street
	}

	nu, ok := users.Update(*um)

	return nu, ok
}

type userDB struct {
	users map[time.Time]userModel
	lock  sync.RWMutex
}

func (udb *userDB) Insert(u userModel) bool {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	for _, v := range udb.users {
		if v.Email == u.Email {
			return false
		}
	}

	udb.users[time.Now()] = u

	return true
}

func (udb *userDB) ReadByEmail(k string) (userModel, bool) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.users {
		if v.Email == k {
			return v, true
		}
	}

	return userModel{}, false
}

func (udb *userDB) ReadById(k string) (userModel, bool) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.users {
		if v.Id == k {
			return v, true
		}
	}

	return userModel{}, false
}

func (udb *userDB) Update(u userModel) (userModel, bool) {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	for _, v := range udb.users {
		if v.Email == u.Email && v.Id != u.Id {
			return userModel{}, false
		}
	}

	for i, v := range udb.users {
		if v.Id == u.Id {
			udb.users[i] = u
			return udb.users[i], true

		}
	}

	return userModel{}, false
}

var users = userDB{users: make(map[time.Time]userModel)}
