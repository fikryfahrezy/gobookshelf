package users

type ForgotPassModel struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}
