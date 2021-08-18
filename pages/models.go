package pages

import (
	"sync"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

type regResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type userSession struct {
	session map[string]string
	lock    sync.RWMutex
}

func (us *userSession) Create(v string) string {
	us.lock.Lock()
	defer us.lock.Unlock()

	k := utils.RandString(15)

	us.session[k] = v

	return k
}

func (us *userSession) Get(k string) string {
	us.lock.RLock()
	defer us.lock.RUnlock()

	return us.session[k]
}

func (us *userSession) Delete(k string) {
	us.lock.Lock()
	defer us.lock.Unlock()

	delete(us.session, k)
}

var userSessions = userSession{session: make(map[string]string)}

const authSessionKey = "auth"
