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
	ProductID uint           `json:"product_id"`
	Product   models.Product `json:"product"`
	Quantity  uint           `json:"quantity"`
	SubTotal  float64        `json:"sub_total"`
}

type ProductIn struct {
	Name          string `json:"name"`
	Quantity      int    `json:"quantity"`
	ProductImages []struct {
		Image string `json:"image"`
	} `json:"product_images"`
	ShipmentTime      string  `json:"shipment_time"`
	Price             float64 `json:"price"`
	SandH             float64 `json:"s_and_h"`
	RSP               int     `json:"rsp"`
	BV                int     `gorm:"default:null" json:"bv"`
	ProductCategoryID uint    `json:"product_category_id"`
	EP                float64 `gorm:"default:null" json:"ep"`
}
