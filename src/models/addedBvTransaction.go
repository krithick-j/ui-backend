package models

import "gorm.io/gorm"

type AddedBv struct {
	gorm.Model
	OrderHeaderID uint `gorm:"foreignKey:OrderHeaderID;references:ID"`
	ReferenceNo   string
	DistribId     string
	Place         string
	Value         float64
}
