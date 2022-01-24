package page

import (
	"sync"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type UserSession struct {
	Session map[string]string
	Lock    sync.RWMutex
}

func (us *UserSession) Create(v string) string {
	us.Lock.Lock()
	defer us.Lock.Unlock()

	k := common.RandString(15)

	us.Session[k] = v

	return k
}

func (us *UserSession) Get(k string) string {
	us.Lock.RLock()
	defer us.Lock.RUnlock()

	return us.Session[k]
}

func (us *UserSession) Delete(k string) {
	us.Lock.Lock()
	defer us.Lock.Unlock()

	delete(us.Session, k)
}
