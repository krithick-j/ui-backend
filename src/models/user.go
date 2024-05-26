package models

import "gorm.io/gorm"

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
	PinOrZipCode    string `json:"pin_or_zip_code"`
	Country         string `json:"country"`
	HomePhoneNo     string `json:"home_phone_no"`
	MobilePhoneNo   string `json:"mobile_phone_no"`
	KYCPhoto        string `json:"kyc_photo"`
	KYCAdhaar       string `json:"kyc_adhaar"`
	KYCPAN          string `json:"kyc_pan"`
	KYCConsentDoc   string `json:"kyc_concesnt_doc"`
	KYCStatus       string `json:"kyc_status" gorm:"default:not_submitted"` //not_submitted | pending | verified
}
