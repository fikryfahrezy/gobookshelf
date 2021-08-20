package galleries

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var pd = "/assets/images"

func saveFileToDir(fn string, f multipart.File) (string, bool) {
	dr, err := os.Getwd()
	if err != nil {
		return err.Error(), false
	}

	l := filepath.Join(dr, pd, fn)
	tf, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err.Error(), false
	}

	defer tf.Close()

	if _, err := io.Copy(tf, f); err != nil {
		return err.Error(), false
	}

	return "", true
}
