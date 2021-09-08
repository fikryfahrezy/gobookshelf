package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/db"
	"github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gobookshelf/geocodings"
	"github.com/fikryfahrezy/gobookshelf/handler"
	"github.com/fikryfahrezy/gobookshelf/pages"
	"github.com/fikryfahrezy/gobookshelf/users"
)

func main() {
	sqliteDb, err := db.InitSqliteDB("data/db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer sqliteDb.Close()

	db.MigrateSqliteDB()

	books.InitDB("data/books.json")

	// Views
	handler.HandlerGET("/", pages.Home)
	handler.HandlerGET("/matrix", pages.Matrix)
	handler.HandlerGET("/register", pages.Register)
	handler.HandlerGET("/logout", pages.Logout)
	handler.HandlerGET("/login", pages.Login)
	handler.HandlerGET("/profile", pages.Profile)
	handler.HandlerGET("/forgotpass", pages.ForgotPass) // handler.HandlerGET("/resetpass", pages.ResetPass) // handler.HandlerGET("/gallery", pages.Gallery)

	// Template Proxy
	handler.HandlerPOST("/registration", pages.Registration)
	handler.HandlerPOST("/loginacc", pages.LoginAcc)
	handler.HandlerPATCH("/updateacc", pages.UpdateAcc)
	handler.HandlerPOST("/oauth", pages.Oauth)

	// Apis
	handler.HandlerPOST("/books", books.Post)
	handler.HandlerGET("/books", books.GetAll)
	handler.HandlerGET("/books/:id", books.GetOne)
	handler.HandlerPUT("/books/:id", books.Put)
	handler.HandlerDELETE("/books/:id", books.Delete)
	handler.HandlerPOST("/userreg", users.Registration)
	handler.HandlerPOST("/userlogin", users.Login)
	handler.HandlerPATCH("/updateuser", users.UpdateProfile)
	handler.HandlerPOST("/forgotpassword", users.ForgotPassword)
	handler.HandlerPATCH("/updatepassword", users.UpdatePassword)
	handler.HandlerGET("/countries", geocodings.GetCountries)
	handler.HandlerGET("/street", geocodings.GetStreet)
	handler.HandlerPOST("/galleries", galleries.Post)
	handler.HandlerGET("/galleries", galleries.Get)

	// Public path
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	handler.InitServer("http://localhost", 8080)
}
