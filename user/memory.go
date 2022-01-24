package user

import (
	"github.com/fikryfahrezy/gobookshelf/common"
	"sync"
)

type Repository struct {
	Users map[string]User
	lock  sync.RWMutex
}

func (r *Repository) GetAll() []User {
	r.lock.Lock()
	defer r.lock.Unlock()

	u := make([]User, len(r.Users))
	i := 0
	for _, user := range r.Users {
		u[i] = user
		i++
	}

	return u
}

func (r *Repository) Insert(u User) User {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Users[common.RandString(9)] = u

	return u
}

func (r *Repository) UpdateById(i string, u User) User {
	r.lock.Lock()
	defer r.lock.Unlock()

	for s, user := range r.Users {
		if user.Id == i {
			r.Users[s] = u
			return u
		}
	}

	return User{}
}
