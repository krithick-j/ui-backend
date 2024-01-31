package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name              string
	Quantity          int
	ProductImages     []ProductImage //group of pictures
	ProductImageID    uint
	ProductImage      ProductImage //main picture
	ShipmentTime      string
	Price             int
	SandH             int
	RSP               int
	BV                int  `gorm:"default:null"`
	AddToCart         bool `gorm:"default:false" `
	ProductCategoryID int
	EP                float32 `gorm:"default:null"`
}
