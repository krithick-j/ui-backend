package models

import "gorm.io/gorm"

//enquiryType := ["Apply ID Card","Bank Account Validation","Commission/ GR Related Enquiry","Ecard/ Evoucher Enquiry","General Enquiry","Grievanvces","GST Enquiry","KYC Submission","","Account/Password/Security Q&W Enquiry",""]

type ContactUs struct {
	gorm.Model
	ContactUsBasicDetails
	EnquiryTypeID uint
	EnquiryField
	ContactQueryFile
}

type ContactUsBasicDetails struct {
	Name          string `json:"name"`
	DistribId     string `json:"distrib_id"`
	Country       string `json:"country"`
	ContactNumber string `json:"contact_number"`
	EmailAddress  string `json:"email_address"`
}

type ContactQueryFile struct {
	AadhaarFront string `json:"aadhaar_front"`
	AadhaarBack  string `json:"aadhaar_back"`
	PanCard      string `json:"pan_card"`
	PassportSize string `json:"passport_size"`
}

type EnquiryType struct {
	gorm.Model
	Name      string `json:"name"`
	AdminName string `json:"admin_name"`
}

type EnquiryField struct {
	YourQuery string `json:"your_query"`
	BankAccountValidation
	RefundRequest
}

type BankAccountValidation struct {
	PanNumber         string `json:"pan_number"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	BankNameID        string `json:"bank_name_id"`
	BankIFSCCode      string `json:"bank_ifsc_code"`
	BankBranch        string `json:"bank_branch"`
}

type BankName struct {
	gorm.Model
	Name string `json:"name"`
}

type RefundRequest struct {
	RefundType     RefundType `json:"refund_type"`
	OrdersHeaderID uint       `json:"orders_header_id"`
}

type RefundType int

const (
	Undefined RefundType = iota
	ProductRefundOnly
	ProductRefundAndNoInterest
)
