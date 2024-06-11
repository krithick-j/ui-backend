package dto

import "ui-back-end/src/models"

// User Input from Client
type UserIn struct {
	Pass string `json:"pass"`
	models.BankDetails
	models.ApplicationInformation
	models.ReferrerInformation
	PreferredPlacementInformationIn
}

type PreferredPlacementInformationIn struct {
	RefPlacementDistribname string `json:"ref_placement_distrib_name"`
	RefPlacementDistribId   string `json:"ref_placement_distrib_id"` //the place or tc in which the user sits in the tree; (PreferredDistribId in models.User)
	RefPlacementPlace       string `json:"ref_placement_place"`      //Preferred Place in models.User
	Side                    string `json:"side"`                     //Preferred Side in models.User
}

type UserOut struct {
	DistribID string `json:"distrib_id"`
}

type AuthOut struct {
	DistribID string `json:"distrib_id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
	KYCStatus string `json:"kyc_status"`
}

type RecursiveUser struct {
	Name           string         `json:"name"`
	TrackingCenter string         `json:"tracking_center"`
	LeftPoint      int            `json:"left_point"`
	RightPoint     int            `json:"right_point"`
	BV             int            `json:"bv"`
	IsActive       bool           `json:"is_active"`
	Left           *RecursiveUser `json:"left"`
	Right          *RecursiveUser `json:"right"`
}

type UserPassIn struct {
	DistribId string `json:"distrib_id"`
	OldPass   string `json:"old_pass"`
	NewPass   string `json:"new_pass"`
}
