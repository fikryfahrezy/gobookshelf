package geocoding

import "net/http"

type CountriesService interface {
	GetCountries(q string) (interface{}, error)
}

type GeoCodeService interface {
	GetGeo(r, s string) (interface{}, error)
}

type Service struct {
	CountriesService
	GeoCodeService
}

type FuncSign func(w http.ResponseWriter, r *http.Request, s *Service)

func (g *Service) Http(f FuncSign) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, g)
	}
}
