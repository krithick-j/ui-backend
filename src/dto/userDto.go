package dto

// User Input from Client
type UserIn struct {
	RefDistID       string `json:"ref_dist_id"`
	Name            string `json:"name"`
	Pass            string `json:"pass"`
	RefCenterCode   string `json:"ref_center_code"`
	Place           string `json:"place"`
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

type UserOut struct {
	DistribID string `json:"ref_dist_id"`
}

type AuthOut struct {
	DistribID string `json:"distrib_id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}
