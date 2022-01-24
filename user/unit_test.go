package user_test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"

	"github.com/fikryfahrezy/gobookshelf/user"

	"github.com/fikryfahrezy/gosrouter"
	_ "modernc.org/sqlite"
)

func TestUsers(t *testing.T) {
	fDb := "./../data/db-test"
	sqldb, err := sql.Open("sqlite", fDb)
	if err != nil {
		log.Fatal(err)
	}

	_, err = sqldb.Exec(common.MigrateSqliteDB())
	if err != nil {
		log.Fatal(err)
	}

	ur := &user.Repository{Users: make(map[string]user.User)}
	fr := &user.ForgotPassRepository{Db: sqldb}
	us := &user.Service{Ur: ur, Fr: fr}
	user.AddRoutes(us)

	cases := []struct {
		testName              string
		init                  func(*http.Request)
		url, method, bodydata string
		expectedCode          int
		expectedResult        int
	}{
		{
			"Registration Success",
			func(r *http.Request) {},
			"/userreg",
			"POST",
			`{"email":"email@email.com","password":"password","name":"name","region":"region","street":"street"}`,
			http.StatusCreated,
			1,
		},
		{
			"Registration Fail, Not Valid Email",
			func(r *http.Request) {},
			"/userreg",
			"POST",
			`{"email":"not-valid-email","password":"password","name":"name","region":"region","street":"street"}`,
			http.StatusUnprocessableEntity,
			1,
		},
		{
			"Login Success",
			func(r *http.Request) {
				u := user.User{
					Email:    "email@email2.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u.Insert(ur)
			},
			"/userlogin",
			"POST",
			`{"email":"email@email2.com","password":"password"}`,
			http.StatusOK,
			2,
		},
		{
			"Login Fail, Password Not Match",
			func(r *http.Request) {
				u := user.User{
					Email:    "email@email3.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u.Insert(ur)
			},
			"/userlogin",
			"POST",
			`{"email":"email@email3.com","password":"not-password"}`,
			http.StatusUnauthorized,
			3,
		},
		{
			"Login Fail, User Not Registered",
			func(r *http.Request) {},
			"/userlogin",
			"POST",
			`{"email":"email@email4.com","password":"password"}`,
			http.StatusUnauthorized,
			3,
		},
		{
			"Login Fail, Not Valid Email",
			func(r *http.Request) {},
			"/userlogin",
			"POST",
			`{"email":"not-valid-email","password":"password"}`,
			http.StatusUnauthorized,
			3,
		},
		{
			"Update Profile Success",
			func(r *http.Request) {
				u := user.User{
					Id:       "1",
					Email:    "email@email4.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u, _ = u.Insert(ur)
				r.Header.Add("authorization", u.Id)
			},
			"/updateprofile",
			"PATCH",
			`{"name":"new name"}`,
			http.StatusOK,
			4,
		},
		{
			"Update Profile Fail, Cookie Not Found",
			func(r *http.Request) {},
			"/updateprofile",
			"PATCH",
			`{"email":"email@email5.com"}`,
			http.StatusUnauthorized,
			4,
		},
		{
			"Request Forgot Password Success",
			func(r *http.Request) {
				u := user.User{
					Email:    "email@email5.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u.Insert(ur)
			},
			"/forgotpassword",
			"POST",
			`{"email":"email@email5.com"}`,
			http.StatusOK,
			5,
		},
		{
			"Update Password Success",
			func(r *http.Request) {
				u := user.User{
					Email:    "email@email6.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u.Insert(ur)

				fp := user.ForgotPass{
					Id:        "1",
					Email:     u.Email,
					Code:      "1",
					IsClaimed: false,
				}
				fp.Insert(fr)
			},
			"/updatepassword",
			"PATCH",
			`{"code":"1", "password":"newpassword"}`,
			http.StatusOK,
			6,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(tr *testing.T) {
			// Use strings.NewReader() because:
			// https://golang.org/pkg/strings/#NewReader
			req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
			if err != nil {
				tr.Fatal(err)
			}

			c.init(req)

			rr := httptest.NewRecorder()
			gosrouter.MakeHandler(rr, req)
			resp := rr.Result()

			if resp.StatusCode != c.expectedCode {
				body, _ := io.ReadAll(resp.Body)
				tr.Log(resp.StatusCode)
				tr.Fatal(string(body))
			}

			if len(ur.Users) != c.expectedResult {
				tr.Fatal(len(ur.Users))
			}
		})
	}

	if err = sqldb.Close(); err != nil {
		return
	}

	if _, err = os.Getwd(); err != nil {
		t.Fatal(err)
	}

	if err = os.Remove(fDb); err != nil {
		t.Fatal(err)
	}
}
