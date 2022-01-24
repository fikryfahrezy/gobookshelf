package page

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"

	"github.com/fikryfahrezy/gosrouter"
)

type UserdataRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func mapUserReqToCmd(u UserdataRequest) UserReqCommand {
	ur := UserReqCommand{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Region:   u.Region,
		Street:   u.Street,
	}
	return ur
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func mapLoginReqToCmd(l LoginRequest) LoginCommand {
	lc := LoginCommand{
		Email:    l.Email,
		Password: l.Password,
	}
	return lc
}

type authTmpl struct {
	OauthURL string
}

func registration(w http.ResponseWriter, r *http.Request, p *Service) {
	var b UserdataRequest
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Message, Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	d, err := p.UserRegistration(mapUserReqToCmd(b))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	if d == "" {
		res := common.Response{Message: "fail", Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	us := p.Session.Create(d)
	res := common.Response{Message: "", Data: d}
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	common.ResJSON(w, http.StatusCreated, res)
}

func loginAcc(w http.ResponseWriter, r *http.Request, p *Service) {
	var b LoginRequest
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Message, Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	d, err := p.UserLogin(mapLoginReqToCmd(b))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusInternalServerError, res)
		return
	}

	if d == "" {
		res := common.Response{Message: "fail", Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	us := p.Session.Create(d)
	res := common.Response{Message: "", Data: d}
	http.SetCookie(w, &http.Cookie{Name: "auth", Value: us, HttpOnly: true, Secure: true, SameSite: 3})
	common.ResJSON(w, http.StatusOK, res)
}

func updateAcc(w http.ResponseWriter, r *http.Request, p *Service) {
	var b UserdataRequest
	errDcd := common.DecodeJSONBody(w, r, &b)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Message, Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	c, err := r.Cookie("auth")
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnauthorized, res)
		return
	}

	uc := p.Session.Get(c.Value)
	d, err := p.UpdateUserAcc(uc, mapUserReqToCmd(b))
	if d == "" {
		res := common.Response{Message: "fail", Data: ""}
		common.ResJSON(w, http.StatusBadRequest, res)
		return
	}

	res := common.Response{Message: "", Data: d}
	common.ResJSON(w, http.StatusOK, res)
}

func oauth(w http.ResponseWriter, r *http.Request, p *Service) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer r.Body.Close()

	http.Redirect(w, r, "/", http.StatusFound)
}

func matrix(w http.ResponseWriter, r *http.Request, p *Service) {
	p.Template.ExecuteTemplate(w, "error.html", nil)
}

func home(w http.ResponseWriter, r *http.Request, p *Service) {
	isLogin := false
	_, err := r.Cookie("auth")
	if err == nil {
		isLogin = true
	}

	q, err := common.ReqQuery(r.URL.String())
	if err != nil {
		http.Redirect(w, r, "/matrix", http.StatusInternalServerError)
		return
	}

	b, err := p.GetBooks(q("name"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	d := struct {
		IsLogin bool
		Books   interface{}
	}{
		isLogin,
		b.Data.Books,
	}
	p.Template.ExecuteTemplate(w, "home.html", d)
}

func register(w http.ResponseWriter, r *http.Request, p *Service) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.Host}
	p.Template.ExecuteTemplate(w, "register.html", a)
}

func logout(w http.ResponseWriter, r *http.Request, p *Service) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	p.Session.Delete(c.Value)
	http.SetCookie(w, &http.Cookie{Name: "auth", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request, p *Service) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a := authTmpl{OauthURL: p.Host}
	p.Template.ExecuteTemplate(w, "login.html", a)
}

func profile(w http.ResponseWriter, r *http.Request, p *Service) {
	c, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	cv := p.Session.Get(c.Value)
	u, err := p.GetUser(cv)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	p.Template.ExecuteTemplate(w, "profile.html", u)
}

func forgotPass(w http.ResponseWriter, r *http.Request, p *Service) {
	_, err := r.Cookie("auth")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.Template.ExecuteTemplate(w, "forgotpass.html", nil)
}

func resetPass(w http.ResponseWriter, r *http.Request, p *Service) {
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

	fpm, err := p.GetForgotPassword(cd)
	if err != nil || fpm.IsClaimed {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p.Template.ExecuteTemplate(w, "resetpass.html", nil)
}

func gallery(w http.ResponseWriter, r *http.Request, p *Service) {
	isLogin := false
	_, err := r.Cookie("auth")
	if err == nil {
		isLogin = true
	}

	im, err := p.GetImages()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	d := struct {
		IsLogin bool
		Images  interface{}
	}{
		isLogin,
		im.Data,
	}
	p.Template.ExecuteTemplate(w, "gallery.html", d)
}

func AddRoutes(r *Service) {
	gosrouter.HandlerPOST("/registration", r.Http(registration))
	gosrouter.HandlerPOST("/loginacc", r.Http(loginAcc))
	gosrouter.HandlerPATCH("/updateacc", r.Http(updateAcc))
	gosrouter.HandlerPOST("/oauth", r.Http(oauth))
	gosrouter.HandlerGET("/", r.Http(home))
	gosrouter.HandlerGET("/matrix", r.Http(matrix))
	gosrouter.HandlerGET("/register", r.Http(register))
	gosrouter.HandlerGET("/logout", r.Http(logout))
	gosrouter.HandlerGET("/login", r.Http(login))
	gosrouter.HandlerGET("/profile", r.Http(profile))
	gosrouter.HandlerGET("/forgotpass", r.Http(forgotPass))
	gosrouter.HandlerGET("/resetpass", r.Http(resetPass))
	gosrouter.HandlerGET("/gallery", r.Http(gallery))
}
