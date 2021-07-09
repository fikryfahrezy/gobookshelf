package books

type bookModelValidator struct {
	Name      string
	Year      int
	Author    string
	Summary   string
	Publisher string
	PageCount int
	ReadPage  int
	Reading   bool
}

func (b *bookModelValidator) Validate() (string, bool) {
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

type getBookQuery struct {
	Name     string
	Reading  string
	Finished string
}
