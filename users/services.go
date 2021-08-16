package users

func createUser(nu userModel) (userModel, bool) {
	cu, ok := nu.Save()

	if !ok {
		return userModel{}, false
	}

	return cu, true
}

func getUser(e string, p string) (userModel, bool) {
	us, ok := users.ReadByEmail(e)

	if !ok || us.Password != p {
		return userModel{}, false
	}

	return us, true
}

func GetUserById(k string) (userModel, bool) {
	us, ok := users.ReadById(k)

	if !ok {
		return userModel{}, false
	}

	return us, true
}

func updateUser(k string, u userModel) (userModel, bool) {
	c, ok := GetUserById(k)

	if !ok {
		return userModel{}, false
	}

	c, ok = c.Update(u)

	return c, ok
}
