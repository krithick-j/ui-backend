package dto

import "time"

type Coupon struct {
	Value    float64 `json:"value"`
	Quantity int64   `json:"quantity"`
}

type ICouponIn struct {
	DistribID string   `json:"distrib_id"`
	TxDetail  string   `json:"tx_detail"`
	Coupons   []Coupon `json:"coupons"`
}

type ValidateICouponIn struct {
	VID string `json:"v_id"`
	Pin string `json:"pin"`
}

type ValidateICouponOut struct {
	Value float64 `json:"value"`
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
	DistribId string    `json:"distrib_id"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
}

type BvHistoryIn struct {
	DistribId   string    `json:"distrib_id"`
	HistoryType string    `json:"history_type"`
	FromDate    time.Time `json:"from_date"`
	ToDate      time.Time `json:"to_date"`
}
