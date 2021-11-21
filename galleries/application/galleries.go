package application

import (
	"fmt"
	"mime/multipart"

	"github.com/fikryfahrezy/gobookshelf/common"
	galleries_common "github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gobookshelf/galleries/domain/galleries"
	galleries_repository "github.com/fikryfahrezy/gobookshelf/galleries/infrastructure/galleries"
)

type GalleryService struct {
	Gr *galleries_repository.ImageRepository
}

func (g GalleryService) InsertImage(f multipart.File, fh multipart.FileHeader) (galleries.GalleryModel, error) {
	alias := common.RandString(8)
	fn := fmt.Sprintf("%s-%s", alias, fh.Filename)

	err := galleries_common.SaveFileToDir(fn, f)
	if err != nil {
		return galleries.GalleryModel{}, err
	}

	im := galleries.GalleryModel{Name: fn}
	res := g.Gr.Insert(im)

	return res, nil
}

func (g GalleryService) GetAllImages() []galleries.GalleryModel {
	return g.Gr.ReadAll()
}

