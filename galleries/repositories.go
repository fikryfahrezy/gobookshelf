package galleries

import (
	"sync"

	"github.com/fikryfahrezy/gobookshelf/common"
)

type imageModel struct {
	Id   string `json:"-"`
	Name string `json:"name"`
}

func (im *imageModel) Save() {
	im.Id = common.RandString(5)

	images.Insert(*im)
}

type imageDB struct {
	images map[string]imageModel
	lock   sync.RWMutex
}

func (idb *imageDB) Insert(im imageModel) string {
	idb.lock.Lock()
	defer idb.lock.Unlock()

	k := common.RandString(7)
	idb.images[k] = im

	return k
}

func (idb *imageDB) ReadAll() []imageModel {
	idb.lock.RLock()
	defer idb.lock.RUnlock()

	ims := make([]imageModel, len(idb.images))
	i := 0

	for _, v := range idb.images {
		ims[i] = v
		i++
	}

	return ims
}

var images = imageDB{images: map[string]imageModel{}}
