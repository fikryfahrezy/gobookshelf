package books

// func Post(w http.ResponseWriter, r *http.Request) {
// 	var b bookReq
// 	errDcd := handler.DecodeJSONBody(w, r, &b)
// 	if errDcd != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: nil}

// 		handler.ResJSON(w, errDcd.Status, res.Response())
// 		return
// 	}

// 	err := b.Validate()
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

// 		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
// 		return
// 	}

// 	nb := bookModel{}
// 	mapBook(&nb, b)

// 	nb = saveBook(nb)
// 	bi := bookIdResponse{nb.Id}
// 	res := handler.CommonResponse{Status: "success", Message: "Book successfully added", Data: bi.Response()}

// 	handler.ResJSON(w, http.StatusCreated, res.Response())
// }

// func GetAll(w http.ResponseWriter, r *http.Request) {
// 	handler.AllowCORS(&w)

// 	q, err := handler.ReqQuery(r.URL.String())
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: make([]interface{}, 0)}

// 		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
// 		return
// 	}

// 	bq := GetBookQuery{q("name"), q("reading"), q("finished")}
// 	b := GetBooks(bq)
// 	bs := booksSerializer{b}
// 	res := handler.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

// 	handler.ResJSON(w, http.StatusOK, res.Response())
// }

// func GetOne(w http.ResponseWriter, r *http.Request) {
// 	p := handler.ReqParams(r.URL.Path)
// 	id := p("id")

// 	if id == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	b, err := getBook(id)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	bs := bookSerializer{b}
// 	res := handler.CommonResponse{Status: "success", Message: "", Data: bs.Response()}

// 	handler.ResJSON(w, http.StatusOK, res.Response())
// }

// func Put(w http.ResponseWriter, r *http.Request) {
// 	var b bookReq
// 	errDcd := handler.DecodeJSONBody(w, r, &b)
// 	if errDcd != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: errDcd.Error(), Data: nil}

// 		handler.ResJSON(w, errDcd.Status, res.Response())
// 		return
// 	}

// 	err := b.Validate()
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

// 		handler.ResJSON(w, http.StatusUnprocessableEntity, res.Response())
// 		return
// 	}

// 	p := handler.ReqParams(r.URL.String())
// 	id := p("id")

// 	if id == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	nb := bookModel{}
// 	mapBook(&nb, b)

// 	nb, err = updateBook(id, nb)

// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	bs := bookSerializer{nb}
// 	res := handler.CommonResponse{Status: "success", Message: "Book successfully updated", Data: bs.Response()}

// 	handler.ResJSON(w, http.StatusOK, res.Response())
// }

// func Delete(w http.ResponseWriter, r *http.Request) {
// 	p := handler.ReqParams(r.URL.Path)
// 	id := p("id")

// 	if id == "" {
// 		res := handler.CommonResponse{Status: "fail", Message: "Not Found", Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	ob, err := deleteBook(id)
// 	if err != nil {
// 		res := handler.CommonResponse{Status: "fail", Message: err.Error(), Data: nil}

// 		handler.ResJSON(w, http.StatusNotFound, res.Response())
// 		return
// 	}

// 	bs := bookSerializer{ob}
// 	res := handler.CommonResponse{Status: "success", Message: "Book successfully deleted", Data: bs.Response()}

// 	handler.ResJSON(w, http.StatusOK, res.Response())
// }
