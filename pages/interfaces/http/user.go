package http

import (
	"html/template"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gobookshelf/pages/application"
	"github.com/fikryfahrezy/gobookshelf/pages/domain/pages"
	"github.com/fikryfahrezy/gobookshelf/users"
	"github.com/fikryfahrezy/gosrouter"
)

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
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
	var b RegistrationRequest
	errDcd := handler.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := handler.CommonResponse{Message: errDcd.Message, Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	cmd := application.RegistrationCommand{
		Email:    b.Email,
		Password: b.Password,
		Name:     b.Name,
		Region:   b.Region,
		Street:   b.Street,
	}

	d, err := p.service.Registration(cmd)
	if err != nil {
		res := handler.CommonResponse{Message: err.Error(), Data: ""}

		handler.ResJSON(w, http.StatusBadRequest, res.Response())
		return
	}

	us := p.session.Create(d)
	res := handler.CommonResponse{Message: "", Data: d}

	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	handler.ResJSON(w, http.StatusCreated, res.Response())
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
	u, err := users.GetUserById(cv)
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

	fpm, err := users.ForgotPasses.ReadByCode(cd)

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

type RegResponseView struct {
	Data string `json:"data"`
}
