package pages

import (
	"html/template"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/users"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func Matrix(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "error.html", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	_, err := r.Cookie(authSessionKey)
	if err == nil {
		isLogin = true
	}

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		http.Redirect(w, r, "/matrix", http.StatusInternalServerError)
		return
	}

	b := books.GetBooks(books.GetBookQuery{Name: q("name")})
	d := struct {
		IsLogin bool
		Books   interface{}
	}{
		isLogin,
		b,
	}

	templates.ExecuteTemplate(w, "home.html", d)
}

func Register(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(authSessionKey)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "register.html", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(authSessionKey)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	userSessions.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{Name: authSessionKey, MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusFound)
}

func Login(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(authSessionKey)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "login.html", nil)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(authSessionKey)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cv := userSessions.Get(c.Value)
	u, ok := users.GetUserById(cv)

	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	templates.ExecuteTemplate(w, "profile.html", u)
}
