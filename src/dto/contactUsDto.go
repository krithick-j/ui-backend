package dto

import (
	"mime/multipart"
	"ui-back-end/src/models"
)

type ContactUsIn struct {
	models.ContactUsBasicDetails
	EnquiryTypeID uint `form:"enquiry_type_id"`
	models.EnquiryField
}

type ContactQueryFileForm struct {
	// AadhaarFront *multipart.File `form:"aadhaar_front"`
	// AadhaarBack  *multipart.File `form:"aadhaar_back"`
	// PanCard      *multipart.File `form:"pan_card"`
	PassportSize *multipart.File `json:"passport_size"`
}

type EnquiryTypeIn struct {
	Name      string           `json:"name"`
	AdminName string           `json:"admin_name"`
	Fields    []EnquiryFieldIn `json:"fields"`
}

type EnquiryFieldIn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
