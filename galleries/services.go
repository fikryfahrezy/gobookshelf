package galleries

import (
	"fmt"
	"mime/multipart"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

func createImage(f multipart.File, fh multipart.FileHeader) (string, bool) {
	alias := utils.RandString(8)
	fn := fmt.Sprintf("%s-%s", alias, fh.Filename)
	msg, ok := saveFileToDir(fn, f)
	if !ok {
		return msg, ok
	}
	im := imageModel{Name: fn}
	im.Save()

	return "", true
}

func GetImages() []imageModel {
	i := images.ReadAll()

	return i
}
