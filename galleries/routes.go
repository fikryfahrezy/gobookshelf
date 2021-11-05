package galleries

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/handler"
)

func Post(w http.ResponseWriter, r *http.Request) {
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

	err = createImage(f, *fh)

	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: ""}

	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func Get(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	i := GetImages()
	ir := imagesResponse{Images: i}
	res := handler.CommonResponse{Message: "", Data: ir.Response()}

	handler.ResJSON(w, http.StatusOK, res.Response())
}
