package dto

// User Input from Client
type UserIn struct {
	Name                  string `json:"name"`
	Pass                  string `json:"pass"`
	RefDistribID          string `json:"ref_distrib_id"`
	Address1              string `json:"address1"`
	Address2              string `json:"address2"`
	TownOrCity            string `json:"town_or_city"`
	District              string `json:"district"`
	StateOrProvince       string `json:"state_or_province"`
	EmailAddress          string `json:"email_address"`
	PinOrZipCode          uint64 `json:"pin_or_zip_code"`
	Country               string `json:"country"`
	HomePhoneNo           string `json:"home_phone_no"`
	MobilePhoneNo         string `json:"mobile_phone_no"`
	RefPlacementDistribId string `json:"ref_placement_distrib_id"` //the place or tc in which the user sits in the tree
	RefPlacementPlace     string `json:"ref_placement_place"`
	Side                  string `json:"side"`
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
