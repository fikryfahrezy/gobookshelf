package gallery_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fikryfahrezy/gobookshelf/gallery"
	"github.com/fikryfahrezy/gosrouter"
)

func TestGalleries(t *testing.T) {
	gli := &gallery.ImageRepository{Images: make(map[string]gallery.Gallery)}
	gls := &gallery.Service{Gr: gli}
	gallery.AddRoutes(gls)

	cases := []struct {
		testName, url, method        string
		expectedCode, expectedResult int
		isUseFormBody                bool
	}{
		{
			"Post Success",
			"/galleries",
			"POST",
			http.StatusCreated,
			1,
			true,
		},
		{
			"Get Success",
			"/galleries",
			"GET",
			http.StatusOK,
			1,
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {

			req, errReq := http.NewRequest(c.method, c.url, nil)

			if c.isUseFormBody {

				gallery.Fd = "../assets/images"

				// REF: POST data using the Content-Type multipart/form-data
				// https://stackoverflow.com/questions/20205796/post-data-using-the-content-type-multipart-form-data
				dr, err := os.Getwd()
				if err != nil {
					panic(err)
				}

				var b bytes.Buffer
				w := multipart.NewWriter(&b)
				fd := filepath.Join(dr, "../fixtures/test-img.jpg")
				f, err := os.Open(fd)
				if err != nil {
					panic(err)
				}

				var ir io.Reader = f
				var fw io.Writer

				if x, ok := ir.(io.Closer); ok {
					defer x.Close()
				}

				// Add image file
				if x, ok := ir.(*os.File); ok {
					ss := strings.Split(x.Name(), "/")
					if fw, err = w.CreateFormFile("image", ss[len(ss)-1]); err != nil {
						panic(err)
					}
				}

				if _, err = io.Copy(fw, ir); err != nil {
					panic(err)
				}

				// Don't forget to close the multipart writer.
				// If you don't close it, your request will be missing the terminating boundary.
				w.Close()

				req, errReq = http.NewRequest(c.method, c.url, &b)

				req.Header.Set("Content-Type", w.FormDataContentType())
			}

			if errReq != nil {
				t.Fatal(errReq)
			}

			rr := httptest.NewRecorder()
			gosrouter.MakeHandler(rr, req)

			if rr.Result().StatusCode != c.expectedCode {
				t.FailNow()
			}

			if len(gli.Images) != c.expectedResult {
				t.FailNow()
			}
		})
	}
}
