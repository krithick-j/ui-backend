package models

import "gorm.io/gorm"

// User struct modified based on Distributor Application Form
type User struct {
	gorm.Model

	DistribID      string  `json:"distrib_id"`
	Pass           string  `json:"pass"`
	CpaPin         string  `json:"cpa_pin"`
	MailingAddress string  `json:"mailing_address"` //aadhar card mail address
	CurrentRank    float64 `json:"current_rank" gorm:"default:1"`
	HighestRank    float64 `json:"highest_rank" gorm:"default:1"`

	ReferrerInformation
	ApplicationInformation
	BankDetails
	PreferredPlacementInformation
	KycDetails
}

type KycDetails struct {
	KYCPhoto                  string `json:"kyc_photo"`
	KYCAdhaar                 string `json:"kyc_adhaar"`
	KYCPAN                    string `json:"kyc_pan"`
	KYCDistribApplicationForm string `json:"kyc_distrib_application_form"`
	KYCAcknowledgementForm    string `json:"kyc_acknowledgement_form"`
	KYCStatus                 string `json:"kyc_status" gorm:"default:not_submitted"` //not_submitted | pending | verified
}

type ApplicationInformation struct {
	Title                   string `json:"title" gorm:"default:'N/A'"` //mr. | mrs. | ms. default is N/A
	Name                    string `json:"name"`
	ChequeName              string `json:"cheque_name"`
	EmailAddress            string `json:"email_address"`
	HomePhoneNo             string `json:"home_phone_no"`
	MobilePhoneNo           string `json:"mobile_phone_no"`
	ValidIdNo               string `json:"valid_id_no"`
	DateOfBirth             string `json:"date_of_birth"`
	MothersMaidenName       string `json:"mothers_maiden_name"`
	BenificiaryName         string `json:"benificiary_name"`
	BeneficiaryRelationship string `json:"benificiary_relationship"`
	AddressDetails
}

type ReferrerInformation struct {
	RefDistribID   string `json:"ref_distrib_id"`
	RefDistribName string `json:"ref_distrib_name"`
}
type BankDetails struct {
	PanCard   string `json:"pan_card"`
	BankName  string `json:"bank_name"`
	BankAccNo string `json:"bank_acc_no"`
	IFSCCode  string `json:"ifsc_code"`
}

type AddressDetails struct {
	Address1        string `json:"address1"` //shipping address
	Address2        string `json:"address2"`
	TownOrCity      string `json:"town_or_city"`
	District        string `json:"district"`
	StateOrProvince string `json:"state_or_province"`
	PinOrZipCode    string `json:"pin_or_zip_code"`
	Country         string `json:"country"`
}

type PreferredPlacementInformation struct {
	PreferredDistribId   string `json:"preferred_distrid_id"`
	PreferredDistribName string `json:"preferred_distrib_name"`
	PreferredPlace       string `json:"preferred_place"`
	PreferredSide        string `json:"preferred_side"`
}
