package users_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	user_service "github.com/fikryfahrezy/gobookshelf/users/application"
	"github.com/fikryfahrezy/gobookshelf/users/infrastructure/forgotpw"
	user_http "github.com/fikryfahrezy/gobookshelf/users/interfaces/http"

	"github.com/fikryfahrezy/gobookshelf/users/domain/users"
	user_repository "github.com/fikryfahrezy/gobookshelf/users/infrastructure/users"

	"github.com/fikryfahrezy/gobookshelf/db"
	"github.com/fikryfahrezy/gosrouter"
)

func TestUsers(t *testing.T) {
	fDb := "./../data/db-test"
	sdb, err := db.InitSqliteTestDB(fDb)
	if err != nil {
		t.FailNow()
	}

	db.MigrateSqliteDB(sdb)

	ur := user_repository.UserRepository{Users: make(map[string]users.UserModel)}
	fr := forgotpw.ForgotPassRepository{Db: sdb}
	us := user_service.UserService{Ur: &ur, Fr: fr}
	usr := user_http.UserRoutes{Us: us}
	user_http.AddRoutes(usr)

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
				u := users.UserModel{
					Email:    "email@email2.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				ur.Insert(u)
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
				u := users.UserModel{
					Email:    "email@email3.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				ur.Insert(u)
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
			http.StatusUnprocessableEntity,
			3,
		},
		{
			"Update Profile Success",
			func(r *http.Request) {
				u := users.UserModel{
					Id:       "1",
					Email:    "email@email4.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				u, _ = ur.Insert(u)
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
				u := users.UserModel{
					Email:    "email@email5.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				ur.Insert(u)
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
				u := users.UserModel{
					Email:    "email@email6.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				ur.Insert(u)

				fp := users.ForgotPassModel{
					Id:        "1",
					Email:     u.Email,
					Code:      "1",
					IsClaimed: false,
				}
				fr.Insert(fp)
			},
			"/updatepassword",
			"PATCH",
			`{"code":"1", "password":"newpassword"}`,
			http.StatusOK,
			6,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			// Use strings.NewReader() because:
			// https://golang.org/pkg/strings/#NewReader
			req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
			if err != nil {
				t.Fatal(err)
			}

			c.init(req)

			rr := httptest.NewRecorder()
			gosrouter.MakeHandler(rr, req)
			resp := rr.Result()

			if resp.StatusCode != c.expectedCode {
				body, _ := io.ReadAll(resp.Body)
				t.Log(resp.StatusCode)
				t.Fatal(string(body))
			}

			if len(ur.Users) != c.expectedResult {
				t.Fatal(len(ur.Users))
			}
		})
	}

	if err = db.RemoveSqliteTestDB(sdb, fDb); err != nil {
		t.Fatal(err)
	}
}
