package dto

type ContactUsIn struct {
	Name          string `json:"name"`
	DistribId     string `json:"distrib_id"`
	Country       string `json:"country"`
	ContactNumber string `json:"contact_number"`
	EmailAddress  string `json:"email_address"`
	EnquiryTypeID uint   `json:"enquiry_type_id"`
	AadhaarFront  string `json:"aadhaar_front"`
	AadhaarBack   string `json:"aadhaar_back"`
	PanCard       string `json:"pan_card"`
	PassportSize  string `json:"passport_size"`
}

type EnquiryTypeIn struct {
	Name      string `json:"name"`
	AdminName string `json:"admin_name"`
	Fields    []EnquiryFieldIn `json:"fields"`
}

type EnquiryFieldIn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
