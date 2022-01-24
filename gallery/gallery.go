package gallery

import "github.com/fikryfahrezy/gobookshelf/common"

type Gallery struct {
	Id   string
	Name string
}

func (g *Gallery) Save(idb *ImageRepository) Gallery {
	k := common.RandString(7)
	g.Id = k

	idb.Insert(*g)

	return *g
}

func (g *Gallery) AllGalleries(idb *ImageRepository) []Gallery {
	gs := idb.ReadAll()
	return gs
}
