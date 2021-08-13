package users

func mapUser(um *userModel, ur regReqValidator) {
	um.Email = ur.Email
	um.Address = ur.Address
	um.Name = ur.Name
	um.Password = ur.Password
}
