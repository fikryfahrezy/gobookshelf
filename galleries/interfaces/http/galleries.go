package http

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/galleries/application"

	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gosrouter"
)

type imageResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func mapImgCmdToResponse(i application.GalleryResCommand) imageResponse {
	ir := imageResponse{
		Id:   i.Id,
		Name: i.Name,
	}
	return ir
}

func mapImgCmdsToResponses(i []application.GalleryResCommand) []imageResponse {
	is := make([]imageResponse, len(i))
	for i, v := range i {
		is[i] = mapImgCmdToResponse(v)
	}
	return is
}

type GalleriesResource struct {
	Service application.GalleryService
}

func (g GalleriesResource) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	f, fh, err := r.FormFile("image")
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	defer f.Close()

	img, err := g.Service.InsertImage(f, *fh)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}
		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: img}
	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func (g GalleriesResource) Get(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	imgs := g.Service.GetAllImages()
	ir := mapImgCmdsToResponses(imgs)
	res := handler.CommonResponse{Message: "", Data: ir}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func AddRoutes(r GalleriesResource) {
	gosrouter.HandlerPOST("/galleries", r.Post)
	gosrouter.HandlerGET("/galleries", r.Get)
}
