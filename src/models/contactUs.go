package models

import "gorm.io/gorm"

//enquiryType := ["Apply ID Card","Bank Account Validation","Commission/ GR Related Enquiry","Ecard/ Evoucher Enquiry","General Enquiry","Grievanvces","GST Enquiry","KYC Submission","","Account/Password/Security Q&W Enquiry",""]

type ContactUs struct {
	Name          string `json:"name"`
	DistribId     string `json:"distrib_id"`
	Country       string `json:"country"`
	ContactNumber string `json:"contact_number"`
	EmailAddress  string `json:"email_address"`
	EnquiryType
	EnquiryTypeID uint   `json:"enquiry_type_id"`
	AadhaarFront  string `json:"aadhaar_front"`
	AadhaarBack   string `json:"aadhaar_back"`
	PanCard       string `json:"pan_card"`
	PassportSize  string `json:"passport_size"`
}

type EnquiryType struct {
	gorm.Model
	Name      string  `json:"name"`
	Fields    []Field `json:"field"`
	AdminName string  `json:"admin_name"`
}

type Field struct {
	gorm.Model
	EnquiryTypeID uint   `json:"enquiry_type_id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Required      bool   `json:"required"`
}

type Enquiry struct {
	gorm.Model
	EnquiryTypeID uint
	EnquiryData   map[string]interface{}
}
