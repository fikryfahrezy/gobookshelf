package gallery

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var Fd = "/assets/images"

func SaveFileToDir(fn string, f multipart.File) error {
	dr, err := os.Getwd()
	if err != nil {
		return err
	}

	l := filepath.Join(dr, Fd, fn)
	tf, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer tf.Close()

	if _, err := io.Copy(tf, f); err != nil {
		return err
	}

	return nil
}
