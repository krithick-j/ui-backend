package models

//enquiryType := ["Apply ID Card","Bank Account Validation","Commission/ GR Related Enquiry","Ecard/ Evoucher Enquiry","General Enquiry","Grievanvces","GST Enquiry","KYC Submission","","Account/Password/Security Q&W Enquiry",""] 

type ContactUs struct{
	DistribId string `json:"distrib_id"`
	Name string `json:"name"`
	Country string `json:"country"`
	ContactNumber string `json:"contact_number"`
	EmailAddress string `json:"email_address"`

}

