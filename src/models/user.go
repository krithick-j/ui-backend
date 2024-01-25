package models

import "gorm.io/gorm"

//	type User struct {
//		ID int `json:"id"`
//		Title string `json:"title"`
//		LastName string `json:"last_name"`
//		FirstName string `json:"first_name"`
//		ChequeName string `json:"cheque_name"`
//		DistributorID string `json:"distributor_id"`
//		Address1 string `json:"address_1"`
//		Address2 string `json:"address_2"`
//		TownOrCity string `json:"town_or_city"`
//		District string `json:"district"`
//		StateOrProvince string `json:"state_or_province"`
//		PostalOrZipCode int `json:"postal_or_zip_code"`
//		HomePhoneNo string `json:"home_phone_no"`
//		MobilePhoneNo string `json:"mobile_phone_no"`
//		EmailAddress string `json:"email_address"`
//		ValidIDType string `json:"valid_id_type"`
//		ValidIDNo string `json:"valid_id_no"`
//		Nationality string `json:"nationality"`
//		Dob string `json:"dob"`
//		NameOfBenificiaryOrNominee string `json:"name_of_benificiary_or_Nominee"`
//		Relationship string `json:"relationship"`
//		MothersFamilyName string `json:"mothers_family_name"`
//		PanCard string `json:"pan_card"`
//	}

type User struct {
	gorm.Model
	DistID    string
	Name      string
	Pass      string
	RefDistID string
	Place     string
	Lside     string
	Rside     string
}
