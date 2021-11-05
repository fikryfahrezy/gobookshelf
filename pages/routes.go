package pages

// func Registration(w http.ResponseWriter, r *http.Request) {
// 	resp, err := http.Post(fmt.Sprintf("%s/userreg", handler.OwnServerUrl), "application/json", r.Body)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	var ur regResp
// 	err = json.NewDecoder(resp.Body).Decode(&ur)

// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	defer resp.Body.Close()

// 	if ur.Data == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: ur.Message, Data: ""}

// 		handler.ResJSON(w, resp.StatusCode, res.Response())
// 		return
// 	}

// 	us := userSessions.Create(ur.Data)
// 	res := handler.CommonResponse{Status: "success", Message: "", Data: ur.Data}

// 	http.SetCookie(w, &http.Cookie{Name: authSessionKey, Value: us, HttpOnly: true, Secure: true, SameSite: 3})
// 	handler.ResJSON(w, http.StatusCreated, res.Response())
// }

// func LoginAcc(w http.ResponseWriter, r *http.Request) {
// 	resp, err := http.Post(fmt.Sprintf("%s/userlogin", handler.OwnServerUrl), "application/json", r.Body)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	var ur regResp
// 	json.NewDecoder(resp.Body).Decode(&ur)

// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	if ur.Data == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: ur.Message, Data: ""}

// 		handler.ResJSON(w, resp.StatusCode, res.Response())
// 		return
// 	}

// 	us := userSessions.Create(ur.Data)
// 	res := handler.CommonResponse{Status: "success", Message: "", Data: ur.Data}

// 	http.SetCookie(w, &http.Cookie{Name: authSessionKey, Value: us, HttpOnly: true, Secure: true, SameSite: 3})
// 	handler.ResJSON(w, http.StatusOK, res)
// }

// func UpdateAcc(w http.ResponseWriter, r *http.Request) {
// 	c, ec := r.Cookie(authSessionKey)

// 	if ec != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: ec.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusUnauthorized, res.Response())
// 		return
// 	}

// 	uc := userSessions.Get(c.Value)
// 	client := &http.Client{}
// 	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/updateuser", handler.OwnServerUrl), r.Body)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	req.Header.Add("authorization", uc)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: ""}

// 		handler.ResJSON(w, http.StatusInternalServerError, res.Response())
// 		return
// 	}

// 	var ur regResp
// 	json.NewDecoder(resp.Body).Decode(&ur)

// 	if ur.Data == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: ur.Message, Data: ""}

// 		handler.ResJSON(w, resp.StatusCode, res.Response())
// 		return
// 	}

// 	res := handler.CommonResponse{Status: "success", Message: "", Data: ur.Data}
// 	handler.ResJSON(w, http.StatusOK, res)
// }

// func Oauth(w http.ResponseWriter, r *http.Request) {
// 	_, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	defer r.Body.Close()

// 	http.Redirect(w, r, "/", http.StatusFound)
// }
