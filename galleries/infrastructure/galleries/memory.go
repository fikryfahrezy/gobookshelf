package galleries

import (
	"sync"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/galleries/domain/galleries"
)

type ImageRepository struct {
	Images map[string]galleries.Gallery
	Lock   sync.RWMutex
}

func (idb *ImageRepository) Insert(im galleries.Gallery) galleries.Gallery {
	idb.Lock.Lock()
	defer idb.Lock.Unlock()

	k := common.RandString(7)
	im.Id = k
	idb.Images[k] = im

	return im
}

func (idb *ImageRepository) ReadAll() []galleries.Gallery {
	idb.Lock.RLock()
	defer idb.Lock.RUnlock()

	ims := make([]galleries.Gallery, len(idb.Images))
	i := 0

	for _, v := range idb.Images {
		ims[i] = v
		i++
	}

	return ims
}
