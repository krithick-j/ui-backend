package dto

type Coupon struct {
	Value     float64 `json:"value"`
	Quantity  int64   `json:"quantity"`
	DateOn    string  `json:"date_on"`
	ExpiresOn string  `json:"expires_on"`
}

type ICouponIn struct {
	DistribID string   `json:"distrib_id"`
	TxDetail  uint64   `json:"tx_detail"`
	Coupons   []Coupon `json:"coupons"`
}

type ValidateICouponIn struct {
	VID string `json:"v_id"`
	Pin string `json:"pin"`
}

type ValidateICouponOut struct {
	Value int `json:"value"`
}
