package models

import "gorm.io/gorm"

type ICouponTransaction struct {
	gorm.Model
	DistribId string  `json:"distrib_id"`
	VID       string  `json:"vid"`
	Value     float64 `json:"value"`
	Reference string  `json:"reference"`
}
