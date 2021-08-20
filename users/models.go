package users

import (
	"sync"
	"time"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

type userModel struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
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

type forgotPassModel struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}

func (fpM *forgotPassModel) Save() {
	fpM.Id = utils.RandString(4)

	ForgotPasses.Insert(*fpM)
}

func (fpM *forgotPassModel) Update(nfpM forgotPassModel) (forgotPassModel, bool) {
	if nfpM.Code != "" {
		fpM.Code = nfpM.Code
	}

	if nfpM.Email != "" {
		fpM.Email = nfpM.Email
	}

	if nfpM.IsClaimed != fpM.IsClaimed {
		fpM.IsClaimed = nfpM.IsClaimed
	}

	nfpM, ok := ForgotPasses.Update(*fpM)

	return nfpM, ok
}

type forgotPassDB struct {
	users map[time.Time]forgotPassModel
	lock  sync.RWMutex
}

func (fpdb *forgotPassDB) Insert(fp forgotPassModel) {
	fpdb.lock.Lock()
	defer fpdb.lock.Unlock()

	fpdb.users[time.Now()] = fp
}

func (fpdb *forgotPassDB) ReadByCode(k string) (forgotPassModel, bool) {
	for _, v := range fpdb.users {
		if v.Code == k && !v.IsClaimed {
			return v, true
		}
	}

	return forgotPassModel{}, false
}

func (fpdb *forgotPassDB) Update(fp forgotPassModel) (forgotPassModel, bool) {
	fpdb.lock.Lock()
	defer fpdb.lock.Unlock()

	for i, v := range fpdb.users {
		if v.Id == fp.Id {
			fpdb.users[i] = fp
			return fpdb.users[i], true
		}
	}

	return forgotPassModel{}, false
}

var ForgotPasses = forgotPassDB{users: make(map[time.Time]forgotPassModel)}
