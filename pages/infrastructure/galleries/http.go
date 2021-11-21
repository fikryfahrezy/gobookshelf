package galleries

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpClient struct {
	Address string
}

func (h HttpClient) GetImages() (interface{}, error) {

	resp, err := http.Get(fmt.Sprintf("%s/galleries", h.Address))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var r interface{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r, nil
}
