package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	DistribID string
	ProductID uint
	Quantity  uint `gorm:"default:1"`
}
