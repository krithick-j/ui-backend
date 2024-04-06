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
	DistribId       string   `json:"distrib_id"`
	Place           string   `json:"place"`
	Coupons         []Coupon `json:"coupons"`
	ICouponTxDetail uint64   `json:"tx_detail"`
}

type TakeChequeOut struct {
	TotalPoints            float32        `json:"total_points"`
	LeftPoints             float32        `json:"left_points"`
	RightPoints            float32        `json:"right_points"`
	TotalCheckoutFrequency int            `json:"total_checkout_frequency"`
	LeftCheckoutFrequency  int            `json:"left_checkout_frequency"`
	RightCheckoutFrequency int            `json:"right_checkout_frequency"`
	RUser                  *RecursiveUser `json:"recursive_user"`
}

type CheckoutIn struct {
	DistribId string `json:"distrib_id"`
}
