package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/gallery"
	"github.com/fikryfahrezy/gobookshelf/geocoding"
	"github.com/fikryfahrezy/gobookshelf/page"

	"github.com/fikryfahrezy/gobookshelf/book"

	"github.com/fikryfahrezy/gobookshelf/user"
	_ "modernc.org/sqlite"

	"github.com/fikryfahrezy/gosrouter"
)

// content holds our static web server content.
//go:embed templates/*
var content embed.FS

var templates = template.Must(template.ParseFS(content, "templates/*"))

func main() {
	sqldb, err := sql.Open("sqlite", "./data/db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqldb.Close()
	_, err = sqldb.Exec(common.MigrateSqliteDB())
	if err != nil {
		log.Fatal(err)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	addr := fmt.Sprintf(":%s", httpPort)
	hostUrl := fmt.Sprintf("http://localhost%s", addr)

	ps := &page.UserSession{Session: map[string]string{}}
	piu := page.UserHttpClient{Address: hostUrl}
	pig := page.GalleryHttpClient{Address: hostUrl}
	pib := page.BookHttpClient{Address: hostUrl}
	pa := &page.Service{
		Host:           hostUrl,
		Template:       templates,
		Session:        ps,
		UserService:    piu,
		GalleryService: pig,
		BookService:    pib,
	}
	page.AddRoutes(pa)

	ur := &user.Repository{Users: make(map[string]user.User)}
	fr := &user.ForgotPassRepository{Db: sqldb}
	us := &user.Service{Ur: ur, Fr: fr}
	user.AddRoutes(us)

	gc := geocoding.CountryHttpClient{Address: "https://restcountries.com"}
	gg := geocoding.GeocodeHttpClient{Address: "https://geocode.xyz"}
	gs := &geocoding.Service{CountriesService: gc, GeoCodeService: gg}
	geocoding.AddRoutes(gs)

	gli := &gallery.ImageRepository{Images: make(map[string]gallery.Gallery)}
	gls := &gallery.Service{Gr: gli}
	gallery.AddRoutes(gls)

	// How to read/write from/to a file using Go
	// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-a-file-using-go
	fd := "data/book.json"
	fl := &book.FileRepository{Filename: fd}
	if _, err = os.Stat(fd); os.IsNotExist(err) {
		fl.WriteFile([]byte("[]"))
	}
	ba := &book.Service{Fr: fl}
	book.AddRoutes(ba)

	// Public path
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	for r := range gosrouter.Routes {
		http.HandleFunc(r, gosrouter.MakeHandler)
	}

	s := &http.Server{
		Addr: addr,
	}

	log.Fatal(s.ListenAndServe())
}
