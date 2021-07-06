package books

import (
	"io"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var b BookModel
	ok := common.DecodeBody(w, r, &b)
	if !ok {
		return
	}
	saveBook(&b)
	io.WriteString(w, "Allo")
}

func Get(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Allo")
}

func Put(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Allo")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Allo")
}
