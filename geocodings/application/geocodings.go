package application

type countriesService interface {
	GetCountries(q string) (interface{}, error)
}

type geoCodeService interface {
	GetGeo(r, s string) (interface{}, error)
}

type GeocodeService struct {
	CountriesService countriesService
	GeoCodeService   geoCodeService
}

func (g GeocodeService) GetCountries(q string) (interface{}, error) {
	return g.CountriesService.GetCountries(q)
}

func (g GeocodeService) GetGeo(r, s string) (interface{}, error) {
	return g.GeoCodeService.GetGeo(r, s)
}
