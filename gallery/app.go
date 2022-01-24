package gallery

import (
	"fmt"
	"github.com/fikryfahrezy/gobookshelf/common"
	"mime/multipart"
	"net/http"
)

type ReqCommand struct {
	Name string
}

type ResCommand struct {
	Id   string
	Name string
}

func mapGalleryReqToEntity(g ReqCommand) Gallery {
	gg := Gallery{
		Name: g.Name,
	}
	return gg
}

func mapGalleryEntityToResCmd(g Gallery) ResCommand {
	rc := ResCommand{
		Id:   g.Id,
		Name: g.Name,
	}
	return rc
}

func mapGalleryEntitiesToResCmds(g []Gallery) []ResCommand {
	cs := make([]ResCommand, len(g))
	for i, v := range g {
		cs[i] = mapGalleryEntityToResCmd(v)
	}

	return cs
}

type Service struct {
	Gr *ImageRepository
}

type FuncSign func(w http.ResponseWriter, r *http.Request, s *Service)

func (g *Service) Http(f FuncSign) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, g)
	}
}

func (g *Service) InsertImage(f multipart.File, fh multipart.FileHeader) (ResCommand, error) {
	alias := common.RandString(8)
	fn := fmt.Sprintf("%s-%s", alias, fh.Filename)

	err := SaveFileToDir(fn, f)
	if err != nil {
		return ResCommand{}, err
	}

	res := g.Gr.Insert(mapGalleryReqToEntity(ReqCommand{Name: fn}))
	rc := mapGalleryEntityToResCmd(res)
	return rc, nil
}

func (g *Service) GetAllImages() []ResCommand {
	r := mapGalleryEntitiesToResCmds(g.Gr.ReadAll())
	return r
}
