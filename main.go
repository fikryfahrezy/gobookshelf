package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fikryfahrezy/gobookshelf/books"
	"github.com/fikryfahrezy/gobookshelf/db"
	pages_app "github.com/fikryfahrezy/gobookshelf/pages/application"
	"github.com/fikryfahrezy/gobookshelf/pages/infrastructure/pages"
	pages_infra "github.com/fikryfahrezy/gobookshelf/pages/infrastructure/users"
	pages_http "github.com/fikryfahrezy/gobookshelf/pages/interfaces/http"
	"github.com/fikryfahrezy/gosrouter"
)

// content holds our static web server content.
//go:embed assets/* templates/*
var content embed.FS

var templates = template.Must(template.ParseFS(content, "templates/*"))

func main() {
	sqliteDb, err := db.InitSqliteDB("data/db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer sqliteDb.Close()

	db.MigrateSqliteDB()

	books.InitDB("data/books.json")

	// Apis
	// gosrouter.HandlerPOST("/books", books.Post)
	// gosrouter.HandlerGET("/books", books.GetAll)
	// gosrouter.HandlerGET("/books/:id", books.GetOne)
	// gosrouter.HandlerPUT("/books/:id", books.Put)
	// gosrouter.HandlerDELETE("/books/:id", books.Delete)
	// gosrouter.HandlerPOST("/userreg", users.Registration)
	// gosrouter.HandlerPOST("/userlogin", users.Login)
	// gosrouter.HandlerPATCH("/updateuser", users.UpdateProfile)
	// gosrouter.HandlerPOST("/forgotpassword", users.ForgotPassword)
	// gosrouter.HandlerPATCH("/updatepassword", users.UpdatePassword)
	// gosrouter.HandlerGET("/countries", geocodings.GetCountries)
	// gosrouter.HandlerGET("/street", geocodings.GetStreet)
	// gosrouter.HandlerPOST("/galleries", galleries.Post)
	// gosrouter.HandlerGET("/galleries", galleries.Get)

	ps := pages.NewUserSession()
	ph := pages_infra.NewHTTPClient("http://localhost:3000")
	pa := pages_app.NewPagesServices(ph)
	pages_http.AddRoutes("http://localhost:3000", pa, ps, templates)

	// Public path
	http.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.FS(content))))

	for r := range gosrouter.Routes {
		http.HandleFunc(r, gosrouter.MakeHandler)
	}

	s := &http.Server{
		Addr: "localhost:3000",
	}

	log.Fatal(s.ListenAndServe())
}
