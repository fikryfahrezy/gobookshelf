package common

import (
	"encoding/json"
	"os"
)

var Filename = "data/books.json"

func saveFile(f string, b []byte) {
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

	if _, err := os.Stat(Filename); os.IsNotExist(err) {

		saveFile(Filename, initdata)

		return
	}
}

func Insert(v interface{}) {
	// open input file
	fi, err := os.Open(Filename)
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

	saveFile(Filename, b)
}
