package books

type BookModel struct {
	Id         string
	Name       string
	Year       int
	Author     string
	Summary    string
	Publisher  string
	PageCount  int
	ReadPage   int
	Finished   bool
	Reading    bool
	IsDeleted  bool
	InsertedAt string
	UpdatedAt  string
}
