package http

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/geocodings/application"
	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gosrouter"
)

type GeocodingsResource struct {
	Service application.GeocodeService
}

func (g GeocodingsResource) GetCountries(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	q, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: make([]interface{}, 0)}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res, err := g.Service.GetCountries(q("name"))
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	handler.ResJSON(w, http.StatusOK, res)
}

func (g GeocodingsResource) GetStreet(w http.ResponseWriter, r *http.Request) {
	handler.AllowCORS(&w)

	q, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	rg, s := q("region"), q("street")
	if rg == "" || s == "" {
		res := handler.CommonResponse{Message: "query needed", Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	res, err := g.Service.GetGeo(rg, s)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: nil}
		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
		return
	}

	handler.ResJSON(w, http.StatusOK, res)
}

func AddRoutes(g GeocodingsResource) {
	gosrouter.HandlerGET("/countries", g.GetCountries)
	gosrouter.HandlerGET("/street", g.GetStreet)
}
