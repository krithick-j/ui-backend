package models

import (
	"time"

	"gorm.io/gorm"
)

type CpaTransaction struct {
	gorm.Model
	DistribId    string    `json:"distrib_id"`
	Reference    string    `json:"reference"`
	ActivateDate time.Time `json:"activate_date"` //activate date in bv tree 15 days
	IsActive     bool      `json:"is_active"`
	Amount       float64   `json:"amount"`
}
