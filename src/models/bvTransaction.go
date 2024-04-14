package models

import (
	"time"

	"gorm.io/gorm"
)

type BvTransaction struct {
	gorm.Model
	DisribId     string    `json:"distrib_id"`
	Place        string    `json:"place"`
	OrderId      string    `json:"order_id"`
	Date         time.Time `json:"date"`
	BvValue      float64   `json:"bv_value"`
	ActivateDate time.Time `json:"activate_date"`
	Side         string    `json:"side"`
	TransType    string    `json:"trans_type"`
}
