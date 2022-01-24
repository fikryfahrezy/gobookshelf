package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GeocodeHttpClient struct {
	Address string
}

func (h GeocodeHttpClient) GetGeo(r, s string) (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/?geoit=json&region=%s&streetname=%s", h.Address, r, s), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
