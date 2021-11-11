package http

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gobookshelf/pages/application"
	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
	"github.com/fikryfahrezy/gosrouter"
)

type UserdataRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type pagesResource struct {
	service  application.PagesService
	session  pages.Session
	template *template.Template
	host     string
}

type authTmpl struct {
	OauthURL string
}

func (p pagesResource) registration(w http.ResponseWriter, r *http.Request) {
	var b UserdataRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	cmd := application.UserCommand(b)
	d, err := p.service.Registration(cmd)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	if d == "" {
		res := handler.CommonResponse{Message: "fail", Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	us := p.session.Create(d)
	res := handler.CommonResponse{Message: "", Data: d}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func (p pagesResource) loginAcc(w http.ResponseWriter, r *http.Request) {
	var b LoginRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	cmd := application.LoginCommand(b)
	d, err := p.service.Login(cmd)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
		return
	}

	if d == "" {
		res := handler.CommonResponse{Message: "fail", Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	us := p.session.Create(d)
	res := handler.CommonResponse{Message: "", Data: d}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	handler.ResJSON(w, http.StatusOK, res)
}

func (p pagesResource) updateAcc(w http.ResponseWriter, r *http.Request) {
	var b UserdataRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	c, err := r.Cookie("auth")
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
		return
	}

	uc := p.session.Get(c.Value)
	cmd := application.UserCommand(b)
	d, err := p.service.UpdateAcc(uc, cmd)

	if d == "" {
		res := handler.CommonResponse{Message: "fail", Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: d}
	handler.ResJSON(w, http.StatusOK, res)
}

func (p pagesResource) oauth(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer r.Body.Close()

	http.Redirect(w, r, "/", http.StatusFound)
}

func (p pagesResource) matrix(w http.ResponseWriter, r *http.Request) {
	p.template.ExecuteTemplate(w, "error.html", nil)
}

func (p pagesResource) home(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	_, err := r.Cookie("auth")
	if err == nil {
		isLogin = true
	}

	q, err := handler.ReqQuery(r.URL.String())
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

	p.template.ExecuteTemplate(w, "home.html", d)
}

func (p pagesResource) register(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.host}

	p.template.ExecuteTemplate(w, "register.html", a)
}

func (p pagesResource) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	p.session.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{Name: "auth", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (p pagesResource) login(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.host}

	p.template.ExecuteTemplate(w, "login.html", a)
}

func (p pagesResource) profile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cv := p.session.Get(c.Value)
	u, err := p.service.GetUser(cv)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	p.template.ExecuteTemplate(w, "profile.html", u)
}

func (p pagesResource) forgotPass(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.template.ExecuteTemplate(w, "forgotpass.html", nil)
}

func (p pagesResource) resetPass(w http.ResponseWriter, r *http.Request) {
	c, err := handler.ReqQuery(r.URL.String())
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cd := c("code")

	if cd == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fpm, err := p.service.GetForgotPassword(cd)

	if err != nil || fpm.IsClaimed {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.template.ExecuteTemplate(w, "resetpass.html", nil)
}

func (p pagesResource) gallery(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	_, err := r.Cookie("auth")
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

	p.template.ExecuteTemplate(w, "gallery.html", d)
}

func AddRoutes(h string, sr application.PagesService, ss pages.Session, t *template.Template) {
	r := pagesResource{host: h, service: sr, session: ss, template: t}
	gosrouter.HandlerPOST("/registration", r.registration)
	gosrouter.HandlerPOST("/loginacc", r.loginAcc)
	gosrouter.HandlerPATCH("/updateacc", r.updateAcc)
	gosrouter.HandlerPOST("/oauth", r.oauth)
	gosrouter.HandlerGET("/", r.home)
	gosrouter.HandlerGET("/matrix", r.matrix)
	gosrouter.HandlerGET("/register", r.register)
	gosrouter.HandlerGET("/logout", r.logout)
	gosrouter.HandlerGET("/login", r.login)
	gosrouter.HandlerGET("/profile", r.profile)
	gosrouter.HandlerGET("/forgotpass", r.forgotPass)
	gosrouter.HandlerGET("/resetpass", r.resetPass)
	gosrouter.HandlerGET("/gallery", r.gallery)
}

