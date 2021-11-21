package http

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/galleries/application"

	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gosrouter"
)

type image struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type imagesResponse struct {
	images []image
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
	ir := imagesResponse{images: make([]image, len(imgs))}
	for i, img := range imgs {
		ir.images[i] = image{
			Id:   img.Id,
			Name: img.Name,
		}
	}

	res := handler.CommonResponse{Message: "", Data: ir}
	handler.ResJSON(w, http.StatusOK, res.Response())
}

func AddRoutes(r GalleriesResource) {
	gosrouter.HandlerPOST("/galleries", r.Post)
	gosrouter.HandlerGET("/galleries", r.Get)
}
