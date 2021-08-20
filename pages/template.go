package pages

import (
	"html/template"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/galleries"
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

	a := authTmplDt{OauthURL: common.OwnServerUrl}

	templates.ExecuteTemplate(w, "register.html", a)
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

	a := authTmplDt{OauthURL: common.OwnServerUrl}

	templates.ExecuteTemplate(w, "login.html", a)
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

func ForgotPass(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(authSessionKey)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "forgotpass.html", nil)
}

func ResetPass(w http.ResponseWriter, r *http.Request) {
	c, err := common.ReqQuery(r.URL.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cd := c("code")

	if cd == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fpm, ok := users.ForgotPasses.ReadByCode(cd)

	if !ok || fpm.IsClaimed {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	templates.ExecuteTemplate(w, "resetpass.html", nil)
}

func Gallery(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	_, err := r.Cookie(authSessionKey)
	if err == nil {
		isLogin = true
	}

	im := galleries.GetImages()
	d := struct {
		IsLogin bool
		Images  interface{}
	}{
		isLogin,
		im,
	}

	templates.ExecuteTemplate(w, "gallery.html", d)
}
