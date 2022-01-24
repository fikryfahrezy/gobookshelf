package gallery

import "sync"

type ImageRepository struct {
	Images map[string]Gallery
	Lock   sync.RWMutex
}

func (idb *ImageRepository) Insert(im Gallery) Gallery {
	idb.Lock.Lock()
	defer idb.Lock.Unlock()

	idb.Images[im.Id] = im

	return im
}

func (idb *ImageRepository) ReadAll() []Gallery {
	idb.Lock.RLock()
	defer idb.Lock.RUnlock()

	ims := make([]Gallery, len(idb.Images))
	i := 0

	for _, v := range idb.Images {
		ims[i] = v
		i++
	}

	return ims
}
