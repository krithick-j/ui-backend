package models

import (
	"gorm.io/gorm"
)

type ICoupon struct {
	gorm.Model
	DistribID string  `json:"distrib_id"`
	DateOn    string  `json:"date_on"`
	TxDetail  uint64  `json:"tx_detail"`
	AdminName string  `json:"admin_name"`
	VID       string  `json:"v_id"`
	Value     float64 `json:"value"`
	ExpiresOn string  `json:"expires_on"`
	Pin       string  `json:"pin"`
}
