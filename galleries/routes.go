package galleries

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024)
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	f, fh, err := r.FormFile("image")
	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	defer f.Close()

	err = createImage(f, *fh)

	if err != nil {
		res := common.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

		common.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	res := common.CommonResponse{Status: "success", Message: "", Data: ""}

	common.ResJSON(w, http.StatusCreated, res.Response())
}

func Get(w http.ResponseWriter, r *http.Request) {
	common.AllowCORS(&w)

	i := GetImages()
	ir := imagesResponse{Images: i}
	res := common.CommonResponse{Status: "success", Message: "", Data: ir.Response()}

	common.ResJSON(w, http.StatusOK, res.Response())
}
