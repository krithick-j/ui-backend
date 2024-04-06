package models

import "gorm.io/gorm"

type ChequeFrequency struct {
	gorm.Model
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	Frequency int    `json:"frequency"`
}
