package users

func mapUser(um *userModel, ur userReq) {
	um.Email = ur.Email
	um.Address = ur.Address
	um.Name = ur.Name
	um.Password = ur.Password
}
