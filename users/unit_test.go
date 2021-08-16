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
		init                  func()
		url, method, bodydata string
		expectedCode          int
		exptedResult          int
		cookieKey             string
		cookieVal             string
	}{
		{
			"Registration Success",
			func() {},
			"/registration",
			"POST",
			`{"email":"email@email.com","password":"password","name":"name","address":"address"}`,
			http.StatusCreated,
			1,
			"c",
			"",
		},
		{
			"Registration Fail, Not Valid Email",
			func() {},
			"/registration",
			"POST",
			`{"email":"not-valid-email","password":"password","name":"name","address":"address"}`,
			http.StatusUnprocessableEntity,
			1,
			"c",
			"",
		},
		{
			"Login Success",
			func() {
				u := userModel{
					Email:    "email@email2.com",
					Password: "password",
					Name:     "Name",
					Address:  "Adress",
				}
				createUser(u)
			},
			"/userlogin",
			"POST",
			`{"email":"email@email2.com","password":"password"}`,
			http.StatusOK,
			2,
			"c",
			"",
		},
		{
			"Login Fail, Password Not Match",
			func() {
				u := userModel{
					Email:    "email@email3.com",
					Password: "password",
					Name:     "Name",
					Address:  "Adress",
				}
				createUser(u)
			},
			"/userlogin",
			"POST",
			`{"email":"email@email3.com","password":"not-password"}`,
			http.StatusUnauthorized,
			3,
			"c",
			"",
		},
		{
			"Login Fail, User Not Registered",
			func() {},
			"/userlogin",
			"POST",
			`{"email":"email@email4.com","password":"password"}`,
			http.StatusUnauthorized,
			3,
			"c",
			"",
		},
		{
			"Login Fail, Not Valid Email",
			func() {},
			"/userlogin",
			"POST",
			`{"email":"not-valid-email","password":"password"}`,
			http.StatusUnprocessableEntity,
			3,
			"c",
			"",
		},
		{
			"Update Profile Success",
			func() {
				u := userModel{
					Id:       "1",
					Email:    "email@email4.com",
					Password: "password",
					Name:     "Name",
					Address:  "Adress",
				}
				users.Insert(u)
				UserSessions.session["1"] = u.Id
			},
			"/updateprofile",
			"PUT",
			`{"name":"new name"}`,
			http.StatusOK,
			4,
			AuthSessionKey,
			"1",
		},
		{
			"Update Profile Fail, Cookie Not Found",
			func() {},
			"/updateprofile",
			"PUT",
			`{"email":"email@email5.com"}`,
			http.StatusUnauthorized,
			4,
			"c",
			"",
		},
	}

	common.HandlerPOST("/registration", Registration)
	common.HandlerPOST("/userlogin", Login)
	common.HandlerPUT("/updateprofile", UpdateProfile)

	for _, c := range cases {
		c.init()

		// Use strings.NewReader() because:
		// https://golang.org/pkg/strings/#NewReader
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
		req.AddCookie(&http.Cookie{Name: c.cookieKey, Value: c.cookieVal})

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		common.MakeHandler(rr, req)

		if rr.Result().StatusCode != c.expectedCode {
			t.Log(rr.Result().StatusCode)
			t.FailNow()
		}

		if len(users.users) != c.exptedResult {
			t.FailNow()
		}
	}
}
