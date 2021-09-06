package galleries

import (
	"fmt"
	"mime/multipart"

	"github.com/fikryfahrezy/gobookshelf/utils"
)

func createImage(f multipart.File, fh multipart.FileHeader) error {
	alias := utils.RandString(8)
	fn := fmt.Sprintf("%s-%s", alias, fh.Filename)
	err := saveFileToDir(fn, f)
	if err != nil {
		return err
	}
	im := imageModel{Name: fn}
	im.Save()

	return nil
}

func GetImages() []imageModel {
	i := images.ReadAll()

	return i
}
