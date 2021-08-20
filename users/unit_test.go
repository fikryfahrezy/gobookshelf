package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func TestHandlers(t *testing.T) {
	users.users = make(map[time.Time]userModel)

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
				u := userModel{
					Email:    "email@email2.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				createUser(u)
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
				u := userModel{
					Email:    "email@email3.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}
				createUser(u)
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
				u := userModel{
					Id:       "1",
					Email:    "email@email4.com",
					Password: "password",
					Name:     "Name",
					Region:   "Region",
					Street:   "Street",
				}

				users.Insert(u)
				r.Header.Add("authorization", "1")
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
	}

	common.HandlerPOST("/userreg", Registration)
	common.HandlerPOST("/userlogin", Login)
	common.HandlerPATCH("/updateprofile", UpdateProfile)

	for _, c := range cases {

		// Use strings.NewReader() because:
		// https://golang.org/pkg/strings/#NewReader
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
		if err != nil {
			t.Fatal(err)
		}

		c.init(req)

		rr := httptest.NewRecorder()
		common.MakeHandler(rr, req)

		if rr.Result().StatusCode != c.expectedCode {
			t.FailNow()
		}

		if len(users.users) != c.expectedResult {
			t.FailNow()
		}
	}
}
