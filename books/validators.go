package books

import "errors"

type bookReq struct {
	Name      string `json:"name"`
	Year      int    `json:"year"`
	Author    string `json:"author"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
	PageCount int    `json:"pageCount"`
	ReadPage  int    `json:"readPage"`
	Reading   bool   `json:"reading"`
}

func (b *bookReq) Validate() error {
	if b.Name == "" {
		return errors.New("name cannot be empty")
	}

	if b.Author == "" {
		return errors.New("author cannot be empty")
	}

	if b.Summary == "" {
		return errors.New("summary cannot be empty")
	}

	if b.Publisher == "" {
		return errors.New("publisher cannot be empty")
	}

	if b.ReadPage > b.PageCount {
		return errors.New("read page cannot be bigger than page count")
	}

	return nil
}

type GetBookQuery struct {
	Name     string
	Reading  string
	Finished string
}
