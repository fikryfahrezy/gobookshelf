package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	user_service "github.com/fikryfahrezy/gobookshelf/users/application"
	"github.com/fikryfahrezy/gobookshelf/users/domain/users"
	"github.com/fikryfahrezy/gobookshelf/users/infrastructure/forgotpw"
	user_infra "github.com/fikryfahrezy/gobookshelf/users/infrastructure/users"
	user_http "github.com/fikryfahrezy/gobookshelf/users/interfaces/http"

	books_app "github.com/fikryfahrezy/gobookshelf/books/application"
	books_infra "github.com/fikryfahrezy/gobookshelf/books/infrastructure/books"
	books_http "github.com/fikryfahrezy/gobookshelf/books/interfaces/http"

	geocodings_app "github.com/fikryfahrezy/gobookshelf/geocodings/application"
	geocodings_infra_countries "github.com/fikryfahrezy/gobookshelf/geocodings/infrastructure/countries"
	geocodings_infra_geocode "github.com/fikryfahrezy/gobookshelf/geocodings/infrastructure/geocode"
	geocodings_http "github.com/fikryfahrezy/gobookshelf/geocodings/interfaces/http"

	galleries_app "github.com/fikryfahrezy/gobookshelf/galleries/application"
	"github.com/fikryfahrezy/gobookshelf/galleries/domain/galleries"
	galleries_infra "github.com/fikryfahrezy/gobookshelf/galleries/infrastructure/galleries"
	galleries_http "github.com/fikryfahrezy/gobookshelf/galleries/interfaces/http"

	"github.com/fikryfahrezy/gobookshelf/db"
	pages_app "github.com/fikryfahrezy/gobookshelf/pages/application"
	pages_infra_books "github.com/fikryfahrezy/gobookshelf/pages/infrastructure/books"
	pages_infra_galleries "github.com/fikryfahrezy/gobookshelf/pages/infrastructure/galleries"
	"github.com/fikryfahrezy/gobookshelf/pages/infrastructure/pages"
	pages_infra_users "github.com/fikryfahrezy/gobookshelf/pages/infrastructure/users"
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
	db.MigrateSqliteDB(sqliteDb)

	ps := &pages.UserSession{Session: map[string]string{}}
	piu := pages_infra_users.HttpClient{Address: "http://localhost:3000"}
	pig := pages_infra_galleries.HttpClient{Address: "http://localhost:3000"}
	pib := pages_infra_books.HttpClient{Address: "http://localhost:3000"}
	pa := pages_app.PagesService{UserService: piu, GalleryService: pig, BookService: pib}
	pr := pages_http.PagesResource{Host: "http://localhost:3000", Service: pa, Session: ps, Template: templates}
	pages_http.AddRoutes(pr)

	ur := user_infra.UserRepository{Users: make(map[string]users.User)}
	fr := forgotpw.ForgotPassRepository{Db: sqliteDb}
	us := user_service.UserService{Ur: &ur, Fr: fr}
	usr := user_http.UserRoutes{Us: us}
	user_http.AddRoutes(usr)

	gc := geocodings_infra_countries.HttpClient{Address: "https://restcountries.com"}
	gg := geocodings_infra_geocode.HttpClient{Address: "https://geocode.xyz"}
	gs := geocodings_app.GeocodeService{CountriesService: gc, GeoCodeService: gg}
	gr := geocodings_http.GeocodingsResource{Service: gs}
	geocodings_http.AddRoutes(gr)

	gli := galleries_infra.ImageRepository{Images: make(map[string]galleries.Gallery)}
	gls := galleries_app.GalleryService{Gr: &gli}
	glr := galleries_http.GalleriesResource{Service: gls}
	galleries_http.AddRoutes(glr)

	bi := books_infra.InitDB("data/books.json")
	ba := books_app.BookService{Fr: bi}
	bh := books_http.BookResource{Service: ba}
	books_http.AddRoutes(bh)

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
