package books

import (
	"net/url"
	"strconv"
)

type BookModelValidator struct {
	Name      string
	Year      int
	Author    string
	Summary   string
	Publisher string
	PageCount int
	ReadPage  int
	Reading   bool
}

func (b *BookModelValidator) Validate() (string, bool) {
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

type (
	QueryVal     struct{}
	GetBookQuery struct {
		Name struct {
			Exist bool
			Val   string
		}
		Reading struct {
			Exist, Val bool
		}
		Finished struct {
			Exist, Val bool
		}
	}
)

func getAllQuery(ur string) (GetBookQuery, error) {
	u, err := url.Parse(ur)
	if err != nil {
		return GetBookQuery{}, err
	}
	q := u.Query()
	var bq GetBookQuery
	if n := q.Get("name"); n != "" {
		bq.Name.Val = n
		bq.Name.Exist = true
	}
	s, err := strconv.ParseInt(q.Get(("reading")), 0, 0)
	if err == nil {
		switch s {
		case 1:
			bq.Reading.Val = true
		default:
			bq.Reading.Val = false
		}
		bq.Reading.Exist = true
	}
	s, err = strconv.ParseInt(q.Get(("finished")), 0, 0)
	if err == nil {
		switch s {
		case 1:
			bq.Finished.Val = true
		default:
			bq.Finished.Val = false
		}
		bq.Finished.Exist = true
	}
	return bq, nil
}
