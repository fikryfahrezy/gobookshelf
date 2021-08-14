package users

import (
	"sync"
	"time"
)

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

func (udb *userDB) Update(u userModel) bool {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	var currU *userModel

	for _, v := range udb.users {
		if v.Email == u.Email {
			currU = &v
			break
		}
	}

	if u.Email != "" {
		currU.Email = u.Email
	}

	if u.Password != "" {
		currU.Password = u.Password
	}

	if u.Name != "" {
		currU.Name = u.Name
	}

	if u.Address != "" {
		currU.Address = u.Address
	}

	return true
}

var users = userDB{users: make(map[time.Time]userModel)}

type userModel struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Address  string `json:"address"`
}

func (um *userModel) Save() {
	users.Insert(*um)
}
