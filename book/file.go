package book

import (
	"encoding/json"
	"os"
)

type FileRepository struct {
	Filename string
}

func (r *FileRepository) WriteFile(b []byte) {
	fo, err := os.Create(r.Filename)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = fo.Write(b); err != nil {
		panic(err)
	}
}

func (r *FileRepository) Insert(v interface{}) {
	fi, err := os.Open(r.Filename)
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
	bs, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	r.WriteFile(bs)
}

func (r *FileRepository) Read(v interface{}) {
	fi, err := os.Open(r.Filename)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err = fi.Close(); err != nil {
			panic(err)
		}
	}()

	json.NewDecoder(fi).Decode(&v)
}

func (r FileRepository) Rewrite(v interface{}) {
	fi, err := os.Open(r.Filename)
	if err != nil {
		panic(err)
	}

	// close fi on exit and check for its returned error
	defer func() {
		if err = fi.Close(); err != nil {
			panic(err)
		}
	}()

	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	r.WriteFile(bs)
}
