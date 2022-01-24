package user

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gosrouter"
)

type userReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type forgotPassReq struct {
	Email string `json:"email"`
}

type resetPassReq struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type userResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Street   string `json:"street"`
}

func mapUserReqToCmd(ur userReq) ReqCommand {
	u := ReqCommand{
		Email:    ur.Email,
		Password: ur.Password,
		Name:     ur.Name,
		Region:   ur.Region,
		Street:   ur.Street,
	}

	return u
}

func mapUserResCmdToRes(us ResCommand) userResponse {
	u := userResponse{
		Id:     us.Id,
		Email:  us.Email,
		Name:   us.Name,
		Region: us.Region,
		Street: us.Street,
	}

	return u
}

type ForgotPwResponse struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}

func mapForgotPwResCmdToRes(f ForgotPwResCommand) ForgotPwResponse {
	r := ForgotPwResponse{
		Id:        f.Id,
		Email:     f.Email,
		Code:      f.Code,
		IsClaimed: f.IsClaimed,
	}

	return r
}

func Registration(w http.ResponseWriter, r *http.Request, s *Service) {
	var u userReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: ""}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	ur, err := s.CreateUser(mapUserReqToCmd(u))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	res := common.Response{Message: "", Data: ur.Id}
	common.ResJSON(w, http.StatusCreated, res)
}

func Login(w http.ResponseWriter, r *http.Request, s *Service) {
	var u loginReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: ""}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	ur, err := s.GetUser(loginCommand{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnauthorized, res)
		return
	}

	ures := mapUserResCmdToRes(ur)
	res := common.Response{Message: "", Data: ures.Id}
	common.ResJSON(w, http.StatusOK, res)
}

func GetProfile(w http.ResponseWriter, r *http.Request, s *Service) {
	c := r.Header.Get("authorization")

	if c == "" {
		res := common.Response{Message: http.StatusText(http.StatusUnauthorized), Data: ""}
		common.ResJSON(w, http.StatusUnauthorized, res)
		return
	}

	ur, err := s.GetUserById(c)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusUnprocessableEntity, res)
		return
	}

	user := mapUserResCmdToRes(ur)
	res := common.Response{Message: "", Data: user}
	common.ResJSON(w, http.StatusOK, res)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request, s *Service) {
	var u userReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: ""}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	c := r.Header.Get("authorization")
	if c == "" {
		res := common.Response{Message: http.StatusText(http.StatusUnauthorized), Data: ""}
		common.ResJSON(w, http.StatusUnauthorized, res)
		return
	}

	ur, err := s.UpdateUser(c, mapUserReqToCmd(u))
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusNotFound, res)
		return
	}

	ures := mapUserResCmdToRes(ur)
	res := common.Response{Message: "", Data: ures.Id}
	common.ResJSON(w, http.StatusOK, res)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request, s *Service) {
	var u forgotPassReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: ""}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	_, err := s.CreateForgotPass(forgotPassCommand{Email: u.Email})
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusInternalServerError, res)
		return
	}

	res := common.Response{Message: "", Data: "Hi"}
	common.ResJSON(w, http.StatusOK, res)
}

func GetForgotPassword(w http.ResponseWriter, r *http.Request, s *Service) {
	p := gosrouter.ReqParams(r.URL.Path)
	c := p("code")

	f, err := s.GetForgotPass(c)
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusInternalServerError, res)
		return
	}

	res := common.Response{Message: "", Data: mapForgotPwResCmdToRes(f)}
	common.ResJSON(w, http.StatusOK, res)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request, s *Service) {
	var u resetPassReq
	errDcd := common.DecodeJSONBody(w, r, &u)
	if errDcd != nil {
		res := common.Response{Message: errDcd.Error(), Data: ""}
		common.ResJSON(w, errDcd.Status, res)
		return
	}

	nfpM, err := s.UpdateForgotPass(resetPassCommand{
		Code:     u.Code,
		Password: u.Password,
	})
	if err != nil {
		res := common.Response{Message: err.Error(), Data: ""}
		common.ResJSON(w, http.StatusInternalServerError, res)
		return
	}

	f := mapForgotPwResCmdToRes(nfpM)
	res := common.Response{Message: "", Data: f.Id}
	common.ResJSON(w, http.StatusOK, res)
}

func AddRoutes(s *Service) {
	gosrouter.HandlerPOST("/userreg", s.Http(Registration))
	gosrouter.HandlerGET("/userprofile", s.Http(GetProfile))
	gosrouter.HandlerPOST("/userlogin", s.Http(Login))
	gosrouter.HandlerPATCH("/updateprofile", s.Http(UpdateProfile))
	gosrouter.HandlerPOST("/forgotpassword", s.Http(ForgotPassword))
	gosrouter.HandlerGET("/forgotpassword/:code", s.Http(GetForgotPassword))
	gosrouter.HandlerPATCH("/updatepassword", s.Http(UpdatePassword))
}
