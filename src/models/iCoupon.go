package models

import (
	"time"

	"gorm.io/gorm"
)

type ICoupon struct {
	gorm.Model
	DistribID string    `json:"distrib_id"`
	DateOn    time.Time `json:"date_on"`
	TxDetail  uint64    `json:"tx_detail"`
	AdminName string    `json:"admin_name"`
	VID       string    `json:"v_id"`
	Value     float64   `json:"value"`
	ExpiresOn time.Time `json:"expires_on"`
	Pin       string    `json:"pin"`
}
