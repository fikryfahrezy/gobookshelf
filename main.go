package main

import (
	"net/http"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/geocodings"
	"github.com/fikryfahrezy/gobookshelf/pages"
	"github.com/fikryfahrezy/gobookshelf/users"
)

func main() {
	books.InitDB()

	// Views
	common.HandlerGET("/", pages.Home)
	common.HandlerGET("/matrix", pages.Matrix)
	common.HandlerGET("/register", pages.Registration)
	common.HandlerGET("/login", pages.Login)
	common.HandlerGET("/profile", pages.Profile)

	// Apis
	common.HandlerPOST("/books", books.Post)
	common.HandlerGET("/books", books.GetAll)
	common.HandlerGET("/books/:id", books.GetOne)
	common.HandlerPUT("/books/:id", books.Put)
	common.HandlerDELETE("/books/:id", books.Delete)
	common.HandlerPOST("/registration", users.Registration)
	common.HandlerPOST("/userlogin", users.Login)
	common.HandlerPATCH("/updateprofile", users.UpdateProfile)
	common.HandlerGET("/logout", users.Logout)
	common.HandlerGET("/countries", geocodings.GetCountries)
	common.HandlerGET("/street", geocodings.GetStreet)

	// Public path
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	common.InitServer(8080)
}
