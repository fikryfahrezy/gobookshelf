package users

type ForgotPass struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}
