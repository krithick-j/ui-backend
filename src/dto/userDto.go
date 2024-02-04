package dto

// User Input from Client
type UserIn struct {
	Name          string `json:"name"`
	Pass          string `json:"pass"`
	RefDistID     string `json:"ref_dist_id"`
	RefCenterCode string `json:"ref_center_code"`
	Place         string `json:"place"`
}

type UserOut struct {
	DistribID string `json:"ref_dist_id"`
}

type AuthOut struct {
	DistribID string `json:"distrib_id"`
	Name      string `json:"name"`
	AuthToken string `json:"auth_token"`
}
