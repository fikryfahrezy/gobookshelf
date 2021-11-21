package galleries_test

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

	galleries_common "github.com/fikryfahrezy/gobookshelf/galleries"
	"github.com/fikryfahrezy/gosrouter"

	galleries_app "github.com/fikryfahrezy/gobookshelf/galleries/application"
	galleries_infra "github.com/fikryfahrezy/gobookshelf/galleries/infrastructure/galleries"
	galleries_http "github.com/fikryfahrezy/gobookshelf/galleries/interfaces/http"

	"github.com/fikryfahrezy/gobookshelf/galleries/domain/galleries"
)

func TestGalleries(t *testing.T) {
	gli := galleries_infra.ImageRepository{Images: make(map[string]galleries.GalleryModel)}
	gls := galleries_app.GalleryService{Gr: &gli}
	glr := galleries_http.GalleriesResource{Service: gls}
	galleries_http.AddRoutes(glr)

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

				galleries_common.Fd = "../assets/images"

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
