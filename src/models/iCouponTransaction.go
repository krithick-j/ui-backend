package models

import "gorm.io/gorm"

type ICouponTransaction struct {
	gorm.Model
	DistribId      string  `json:"distrib_id"`
	VID            string  `json:"vid"`
	Pin            string  `json:"pin"`
	TotalValue     float64 `json:"total_value"`
	AmountDetected float64 `json:"amount_detected"`
	Balance        float64 `json:"balance"`
	OrderId        string
}
