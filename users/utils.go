package users

func mapUser(um *userModel, ur userReq) {
	um.Email = ur.Email
	um.Region = ur.Region
	um.Street = ur.Street
	um.Name = ur.Name
	um.Password = ur.Password
}
