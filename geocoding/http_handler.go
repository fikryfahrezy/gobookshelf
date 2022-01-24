package geocoding

import (
	"github.com/fikryfahrezy/gobookshelf/common"
	"net/http"

	"github.com/fikryfahrezy/gosrouter"
)

func GetCountries(w http.ResponseWriter, r *http.Request, g *Service) {
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.Response{Message: err.Error(), Data: make([]interface{}, 0)}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := g.GetCountries(q("name"))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	common.ResJSON(w, http.StatusOK, res)
}

func GetStreet(w http.ResponseWriter, r *http.Request, g *Service) {
	common.AllowCORS(&w)

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	rg, s := q("region"), q("street")
	if rg == "" || s == "" {
		res := common.Response{Message: "query needed", Data: nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	res, err := g.GetGeo(rg, s)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: nil}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	common.ResJSON(w, http.StatusOK, res)
}

func AddRoutes(g *Service) {
	gosrouter.HandlerGET("/countries", g.Http(GetCountries))
	gosrouter.HandlerGET("/street", g.Http(GetStreet))
}
