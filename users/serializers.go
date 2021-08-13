package users

type commonResponse struct{}

func (c *commonResponse) Response() *commonResponse {
	return c
}
