package http

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

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

type PagesResource struct {
	Service  application.PagesService
	Session  pages.Session
	Template *template.Template
	Host     string
}

type authTmpl struct {
	OauthURL string
}

func (p PagesResource) registration(w http.ResponseWriter, r *http.Request) {
	var b UserdataRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	cmd := application.UserCommand(b)
	d, err := p.Service.Registration(cmd)
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

	us := p.Session.Create(d)
	res := handler.CommonResponse{Message: "", Data: d}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	handler.ResJSON(w, http.StatusCreated, res.Response())
}

func (p PagesResource) loginAcc(w http.ResponseWriter, r *http.Request) {
	var b LoginRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	cmd := application.LoginCommand(b)
	d, err := p.Service.Login(cmd)
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

	us := p.Session.Create(d)
	res := handler.CommonResponse{Message: "", Data: d}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	handler.ResJSON(w, http.StatusOK, res)
}

func (p PagesResource) updateAcc(w http.ResponseWriter, r *http.Request) {
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

	uc := p.Session.Get(c.Value)
	cmd := application.UserCommand(b)
	d, err := p.Service.UpdateAcc(uc, cmd)

	if d == "" {
		res := handler.CommonResponse{Message: "fail", Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	res := handler.CommonResponse{Message: "", Data: d}
	handler.ResJSON(w, http.StatusOK, res)
}

func (p PagesResource) oauth(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer r.Body.Close()

	http.Redirect(w, r, "/", http.StatusFound)
}

func (p PagesResource) matrix(w http.ResponseWriter, r *http.Request) {
	p.Template.ExecuteTemplate(w, "error.html", nil)
}

func (p PagesResource) home(w http.ResponseWriter, r *http.Request) {
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

	b, err := p.Service.GetBooks(q("name"))
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	d := struct {
		IsLogin bool
		Books   interface{}
	}{
		isLogin,
		b,
	}

	p.Template.ExecuteTemplate(w, "home.html", d)
}

func (p PagesResource) register(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.Host}

	p.Template.ExecuteTemplate(w, "register.html", a)
}

func (p PagesResource) logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	p.Session.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{Name: "auth", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (p PagesResource) login(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.Host}

	p.Template.ExecuteTemplate(w, "login.html", a)
}

func (p PagesResource) profile(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cv := p.Session.Get(c.Value)
	u, err := p.Service.GetUser(cv)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	p.Template.ExecuteTemplate(w, "profile.html", u)
}

func (p PagesResource) forgotPass(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.Template.ExecuteTemplate(w, "forgotpass.html", nil)
}

func (p PagesResource) resetPass(w http.ResponseWriter, r *http.Request) {
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

	fpm, err := p.Service.GetForgotPassword(cd)

	if err != nil || fpm.IsClaimed {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.Template.ExecuteTemplate(w, "resetpass.html", nil)
}

func (p PagesResource) gallery(w http.ResponseWriter, r *http.Request) {
	isLogin := false
	_, err := r.Cookie("auth")
	if err == nil {
		isLogin = true
	}

	im, err := p.Service.GalleryService.GetImages()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	d := struct {
		IsLogin bool
		Images  interface{}
	}{
		isLogin,
		im,
	}

	p.Template.ExecuteTemplate(w, "gallery.html", d)
}

func AddRoutes(r PagesResource) {
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
