package users

func createUser(nu userModel) userModel {
	nu.Save()

	return nu
}
