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
	DistribID       string `json:"distrib_id"`
	RefDistribID    string `json:"ref_distrib_id"`
	Name            string `json:"name"`
	Pass            string `json:"pass"`
	Address1        string `json:"address1"`
	Address2        string `json:"address2"`
	TownOrCity      string `json:"town_or_city"`
	District        string `json:"district"`
	StateOrProvince string `json:"state_or_province"`
	EmailAddress    string `json:"email_address"`
	PinOrZipCode    uint64 `json:"pin_or_zip_code"`
	Country         string `json:"country"`
	HomePhoneNo     uint64 `json:"home_phone_no"`
	MobilePhoneNo   uint64 `json:"mobile_phone_no"`
}

type TrackingCenter struct {
	gorm.Model
	DistribID      string `json:"distrib_id"`
	Name           string `json:"name"`
	RefDistribID   string `json:"ref_distrib_id"`
	LeftDistribID  string `json:"left_distrib_id"`
	LeftPlace      string `json:"left_place"`
	RightDistribID string `json:"right_distrib_id"`
	RightPlace     string `json:"right_place"`
	CenterCode     string `json:"center_code"`
}
