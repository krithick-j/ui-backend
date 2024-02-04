package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Quantity  uint
}
