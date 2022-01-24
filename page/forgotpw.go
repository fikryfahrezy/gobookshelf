package page

type ForgotPass struct {
	Id        int
	Email     string
	Code      string
	IsClaimed bool
}
