package books

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

// Make variable exported so can be changed in testing
var filename = "data/books.json"

func writeFile(f string, b []byte) {
	fo, err := os.Create(f)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := fo.Write(b); err != nil {
		panic(err)
	}
}

// How to read/write from/to a file using Go
// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
func InitDB() {
	// open input file
	initdata := []byte("[]")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		writeFile(filename, initdata)
	}
}

func insert(v interface{}) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var d []interface{}
	json.NewDecoder(fi).Decode(&d)

	d = append(d, v)
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	writeFile(filename, b)
}

func read(v interface{}) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	json.NewDecoder(fi).Decode(&v)
}

func Update(v interface{}) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	writeFile(filename, b)
}

type bookModel struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Year       int    `json:"year"`
	Author     string `json:"author"`
	Summary    string `json:"summary"`
	Publisher  string `json:"publisher"`
	PageCount  int    `json:"pageCount"`
	ReadPage   int    `json:"readPage"`
	Finished   bool   `json:"finished"`
	Reading    bool   `json:"reading"`
	IsDeleted  bool   `json:"-"`
	InsertedAt string `json:"insertedAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func (b *bookModel) Save() {
	t := time.Now().UTC().String()
	b.Id = utils.RandString(5)
	b.InsertedAt = t
	b.UpdatedAt = t

	insert(b)
}

func (b *bookModel) Update(nb bookModel) {
	var ob bookModel
	bs := GetAllBooks()
	ci := -1

	for i, v := range bs {
		if v.Id == b.Id {
			ci = i
			ob = v
			break
		}
	}

	if ci >= 0 {
		ob.Name = nb.Name
		ob.Year = nb.Year
		ob.Author = nb.Author
		ob.Summary = nb.Summary
		ob.Publisher = nb.Publisher
		ob.PageCount = nb.PageCount
		ob.ReadPage = nb.ReadPage
		ob.Reading = nb.Reading
		ob.Finished = nb.ReadPage == nb.PageCount
		ob.UpdatedAt = time.Now().UTC().String()

		bs[ci] = ob

		Update(bs)
	}
}

func (b *bookModel) Delete() {
	bs := GetAllBooks()

	for i, v := range bs {
		if v.Id == b.Id {
			bs[i].IsDeleted = true
		}
	}

	Update(bs)
}

func GetAllBooks() []bookModel {
	var b []bookModel
	read(&b)

	return b
}

func GetSelectedBooks(q GetBookQuery) []bookModel {
	b := []bookModel{}
	read(&b)

	var nb []bookModel
	n, f, d := q.Name, q.Finished, q.Reading

	for _, v := range b {
		if pf, err := strconv.ParseBool(f); err == nil && pf == v.Finished {
			nb = append(nb, v)
		}
		if pd, err := strconv.ParseBool(d); err == nil && pd == v.Reading {
			nb = append(nb, v)
		}
		if n != "" && strings.Contains(strings.ToLower(v.Name), strings.ToLower(n)) {
			nb = append(nb, v)
		}
	}

	if nb == nil {
		return b
	}

	return nb
}
