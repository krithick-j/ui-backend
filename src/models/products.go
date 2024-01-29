package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name              string
	Quantity          int
	ProductImages     []ProductImage
	ShipmentTime      string
	Price             int
	SandH             int
	RSP               int
	BV                int
	AddToCart         bool
	ProductCategoryID int
	ProductCategory   ProductCategory
}
