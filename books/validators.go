package books

type bookReqValidator struct {
	Name      string `json:"name"`
	Year      int    `json:"year"`
	Author    string `json:"author"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
	PageCount int    `json:"pageCount"`
	ReadPage  int    `json:"readPage"`
	Reading   bool   `json:"reading"`
}

func (b *bookReqValidator) Validate() (string, bool) {
	if b.Name == "" {
		return "Name cannot be empty", false
	}

	if b.Author == "" {
		return "Author cannot be empty", false
	}

	if b.Summary == "" {
		return "Summary cannot be empty", false
	}

	if b.Publisher == "" {
		return "Publisher cannot be empty", false
	}

	if b.ReadPage > b.PageCount {
		return "Read page cannot be bigger than page count", false
	}

	return "", true
}

type GetBookQuery struct {
	Name     string
	Reading  string
	Finished string
}
