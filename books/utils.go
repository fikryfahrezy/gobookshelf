package books

func mapBook(ob *bookModel, nb bookReqValidator) {
	ob.Name = nb.Name
	ob.Year = nb.Year
	ob.Author = nb.Author
	ob.Summary = nb.Summary
	ob.Publisher = nb.Publisher
	ob.PageCount = nb.PageCount
	ob.ReadPage = nb.ReadPage
	ob.Reading = nb.Reading
	ob.Finished = nb.ReadPage == nb.PageCount
}
