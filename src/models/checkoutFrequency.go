package models

import "gorm.io/gorm"

type CheckoutFrequency struct {
	gorm.Model
	DistribId string `json:"distrib_id"`
	Frequency int    `json:"frequency"`
}
