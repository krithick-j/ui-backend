package dto

import "time"

type Coupon struct {
	Value    float64 `json:"value"`
	Quantity int64   `json:"quantity"`
}

type ICouponIn struct {
	DistribID string   `json:"distrib_id"`
	Coupons   []Coupon `json:"coupons"`
}

type ValidateICouponIn struct {
	VID string `json:"v_id"`
	Pin string `json:"pin"`
}

type ValidateICouponOut struct {
	TotalValue float64 `json:"total_value"`
	Value      float64 `json:"value"`
}

type SendCoupon struct {
	VID       string
	Pin       string
	Value     float64
	DateOn    time.Time
	ExpiresOn time.Time
	Active    bool
}

type ICouponHistoryIn struct {
	DistribId string `json:"distrib_id"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
}

type BvHistoryIn struct {
	DistribId string    `json:"distrib_id"`
	TransType string    `json:"trans_type"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
}

type GetICouponOut struct {
	DistribID      string  `json:"distrib_id"`
	DateOn         string  `json:"date_on"`
	Reference      string  `json:"reference"`
	AdminName      string  `json:"admin_name"`
	VID            string  `json:"v_id"`
	TotalValue     float64 `json:"total_value"`
	RemainingValue float64 `json:"remaining_value"`
	ExpiresOn      string  `json:"expires_on"`
	Pin            string  `json:"pin"`
	Active         bool    `json:"active" default:"true"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	ID             uint    `json:"id"`
}
