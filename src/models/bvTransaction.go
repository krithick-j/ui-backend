package models

import (
	"time"

	"gorm.io/gorm"
)

type BvTransaction struct {
	gorm.Model
	DistribId          string    `json:"distrib_id"`
	Place              string    `json:"place"`
	OrderId            string    `json:"order_id"`
	Date               time.Time `json:"date"`
	BvValue            float64   `json:"bv_value"`
	ActivateDate       time.Time `json:"activate_date"` //activate date in bv tree 15 days
	IsActive           int       `json:"is_active"`
	Side               string    `json:"side"`
	TransType          string    `json:"trans_type"`
}
