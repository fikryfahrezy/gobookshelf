package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/common"
)

func TestHandlers(t *testing.T) {
	cases := []struct {
		testName              string
		init                  func()
		url, method, bodydata string
		expectedCode          int
	}{
		{
			"Registration Success",
			func() {},
			"/users",
			"POST",
			`{"email":"email@email.com","password":"password","name":"name","address":"address"}`,
			http.StatusCreated,
		},
	}

	common.HandlerPOST("/users", Registration)

	for _, c := range cases {
		c.init()

		// Use strings.NewReader() because:
		// https://golang.org/pkg/strings/#NewReader
		req, err := http.NewRequest(c.method, c.url, strings.NewReader(c.bodydata))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		common.MakeHandler(rr, req)

		if rr.Result().StatusCode != c.expectedCode {
			t.FailNow()
		}
	}
}
