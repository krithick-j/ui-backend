pack
ge models

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
	Price             float64
	SandH             float64
	RSP               int
	BV                int `gorm:"default:null"`
	ProductCategoryID int
	EP                float64 `gorm:"default:null"`
}
