package api

type PdfResource struct {
	Url string
	Meta Metadata
}

type Metadata struct {
	Title string
	Author string
}