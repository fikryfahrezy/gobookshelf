package application

import (
	"fmt"
	"mime/multipart"

	"github.com/fikryfahrezy/gobookshelf/common"
	galleries_common "github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gobookshelf/galleries/domain/galleries"
	galleries_repository "github.com/fikryfahrezy/gobookshelf/galleries/infrastructure/galleries"
)

type GalleryReqCommand struct {
	Name string
}

type GalleryResCommand struct {
	Id   string
	Name string
}

func mapGalleryReqToEntity(g GalleryReqCommand) galleries.Gallery {
	gg := galleries.Gallery{
		Name: g.Name,
	}
	return gg
}

func mapGalleryEntityToResCmd(g galleries.Gallery) GalleryResCommand {
	rc := GalleryResCommand{
		Id:   g.Id,
		Name: g.Name,
	}
	return rc
}

func mapGalleryEntitiesToResCmds(g []galleries.Gallery) []GalleryResCommand {
	cs := make([]GalleryResCommand, len(g))
	for i, v := range g {
		cs[i] = mapGalleryEntityToResCmd(v)
	}

	return cs
}

type GalleryService struct {
	Gr *galleries_repository.ImageRepository
}

func (g GalleryService) InsertImage(f multipart.File, fh multipart.FileHeader) (GalleryResCommand, error) {
	alias := common.RandString(8)
	fn := fmt.Sprintf("%s-%s", alias, fh.Filename)

	err := galleries_common.SaveFileToDir(fn, f)
	if err != nil {
		return GalleryResCommand{}, err
	}

	res := g.Gr.Insert(mapGalleryReqToEntity(GalleryReqCommand{Name: fn}))
	rc := mapGalleryEntityToResCmd(res)
	return rc, nil
}

func (g GalleryService) GetAllImages() []GalleryResCommand {
	r := mapGalleryEntitiesToResCmds(g.Gr.ReadAll())
	return r
}
