package dto

type Tc struct {
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	LPoint    int    `json:"l_point"`
	RPoint    int    `json:"r_point"`
	BvPoint   int    `json:"bv_point"`
}
type ChequeAvailableIn struct {
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	LTc       Tc     `json:"l_tc"`
	RTc       Tc     `json:"r_tc"`
}

type TakeChequeIn struct {
	DistribId string   `json:"distrib_id"`
	Place     string   `json:"place"`
	Coupons   []Coupon `json:"coupons"`
}

type TakeChequeOut struct {
	TotalBalance           float64          `json:"total_balance"`
	TotalAvailableBalance  float64          `json:"total_available_balance"`
	PlacePointsArr         []PlacePointsArr `json:"place_points"`
	DirectCommissionPoints float64          `json:"direct_commission_points"`
}

type PlacePointsArr struct {
	Place string  `json:"place,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type CheckoutIn struct {
	DistribId string `json:"distrib_id"`
}

type CheckoutFrequency struct {
	TotalCheckoutFrequency  int `json:"total_checkout_frequency"`
	ParentCheckoutFrequency int `json:"parent_checkout_frequency"`
	LeftCheckoutFrequency   int `json:"left_checkout_frequency"`
	RightCheckoutFrequency  int `json:"right_checkout_frequency"`
}
