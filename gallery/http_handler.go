package gallery

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"

	"github.com/fikryfahrezy/gosrouter"
)

type imageResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func mapImgCmdToResponse(i ResCommand) imageResponse {
	ir := imageResponse{
		Id:   i.Id,
		Name: i.Name,
	}
	return ir
}

func mapImgCmdsToResponses(i []ResCommand) []imageResponse {
	is := make([]imageResponse, len(i))
	for i, v := range i {
		is[i] = mapImgCmdToResponse(v)
	}
	return is
}

func Post(w http.ResponseWriter, r *http.Request, g *Service) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	f, fh, err := r.FormFile("image")
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	defer f.Close()

	img, err := g.InsertImage(f, *fh)

	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusInternalServerError, res)
		return
	}

	res := common.Response{Message: "", Data: img}
	common.ResJSON(w, http.StatusCreated, res)
}

func Get(w http.ResponseWriter, r *http.Request, g *Service) {
	common.AllowCORS(&w)

	imgs := g.GetAllImages()
	ir := mapImgCmdsToResponses(imgs)
	res := common.Response{Message: "", Data: ir}
	common.ResJSON(w, http.StatusOK, res)
}

func AddRoutes(g *Service) {
	gosrouter.HandlerPOST("/galleries", g.Http(Post))
	gosrouter.HandlerGET("/galleries", g.Http(Get))
}
