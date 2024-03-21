package dto

import "mime/multipart"

type ContactUsIn struct {
	RequestType  string `json:"request_type"`
	Query        string `json:"query"`
	AadharPdf    *multipart.FileHeader
	PanCard      *multipart.FileHeader
	PassportSize *multipart.FileHeader
}
