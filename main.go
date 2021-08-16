package main

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/data"
	"github.com/fikryfahrezy/gobookshelf/pages"
	"github.com/fikryfahrezy/gobookshelf/users"
)

func rootPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages", http.StatusFound)
}

func main() {
	data.InitDB()
	common.HandlerGET("/", rootPage)
	common.HandlerGET("/home", pages.Home)
	common.HandlerGET("/matrix", pages.Matrix)
	common.HandlerGET("/register", pages.Registration)
	common.HandlerGET("/login", pages.Login)
	common.HandlerGET("/profile", pages.Profile)
	common.HandlerPOST("/books", books.Post)
	common.HandlerGET("/books", books.GetAll)
	common.HandlerGET("/books/:id", books.GetOne)
	common.HandlerPUT("/books/:id", books.Put)
	common.HandlerDELETE("/books/:id", books.Delete)
	common.HandlerPOST("/registration", users.Registration)
	common.HandlerPOST("/userlogin", users.Login)
	common.HandlerPUT("/updateprofile", users.UpdateProfile)
	common.HandlerGET("/logout", users.Logout)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	common.InitServer(8080)
}
