package page

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GalleryHttpClient struct {
	Address string
}

type image struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ImageClientRes struct {
	Data []image `json:"data"`
}

func (h GalleryHttpClient) GetImages() (ImageClientRes, error) {

	resp, err := http.Get(fmt.Sprintf("%s/galleries", h.Address))
	if err != nil {
		return ImageClientRes{}, err
	}

	defer resp.Body.Close()

	var r ImageClientRes
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return ImageClientRes{}, err
	}

	return r, nil
}
