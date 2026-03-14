package api

type PdfResource struct {
	Pdf string
	Meta Metadata
}

type Metadata struct {
	Title string
	Author string
}