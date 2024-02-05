package dto

import "ui-back-end/src/models"

type ProductCart struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type CartItemIn struct {
	DistribID string        `json:"distrib_id"`
	Items     []ProductCart `json:"items"`
}

type ProductsOut struct {
	ProductID uint
	Product   models.Product
	Quantity  uint
	SubTotal  float64
}

type ProductIn struct {
	Name              string
	Quantity          int
	ProductImageID    uint
	ShipmentTime      string
	Price             float64
	SandH             float64
	RSP               int
	BV                int `gorm:"default:null"`
	ProductCategoryID int
	EP                float64 `gorm:"default:null"`
}
