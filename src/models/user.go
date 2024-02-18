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
	PinOrZipCode    uint64 `json:"pin_or_zip_code"`
	Country         string `json:"country"`
	HomePhoneNo     string `json:"home_phone_no"`
	MobilePhoneNo   string `json:"mobile_phone_no"`
}

type TrackingCenter struct {
	gorm.Model
	DistribID      string `json:"distrib_id"`
	Name           string `json:"name"`
	PDistribId     string `json:"p_distrib_id"`
	PPlace         string `json:"p_place"`
	LeftDistribID  string `json:"left_distrib_id"`
	LeftPlace      string `json:"left_place"`
	RightDistribID string `json:"right_distrib_id"`
	RightPlace     string `json:"right_place"`
	Place          string `json:"place"`
	Bv             int    `json:"bv"`
}
