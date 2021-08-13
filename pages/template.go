package pages

import (
	"html/template"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func Page(w http.ResponseWriter, r *http.Request) {
	q, err := common.ReqQuery(r.URL.String())

	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	b := books.GetBooks(books.GetBookQuery{Name: q("name")})
	d := struct{ Books interface{} }{b}

	templates.ExecuteTemplate(w, "index.html", d)
}

func Matrix(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "error.html", nil)
}
