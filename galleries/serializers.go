package galleries

type imagesResponse struct {
	Images []imageModel
}

func (ir *imagesResponse) Response() []imageModel {
	return ir.Images
}
