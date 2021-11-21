package countries

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpClient struct {
	Address string
}

func (h HttpClient) GetCountries(q string) (interface{}, error) {
	var req *http.Request
	client := &http.Client{}

	var err error
	if q == "" {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/v2/all", h.Address), nil)
	} else {
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/v2/name/%s", h.Address, q), nil)
	}

	if err != nil {
		return nil, err
	}

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
