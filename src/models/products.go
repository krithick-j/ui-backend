package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name              string         `json:"name"`
	AdminName         string         `json:"admin_name"`
	Quantity          int            `json:"quantity"`
	ShipmentTime      string         `json:"shipment_time"`
	Price             float64        `json:"price"`
	SandH             float64        `json:"s_and_h"`
	RSP               int            `json:"rsp"`
	BV                int            `gorm:"default:null" json:"bv"`
	ProductCategoryID uint           `json:"product_category_id"`
	EP                float64        `gorm:"default:null" json:"ep"`
	ProductType       string         `gorm:"default:not null" json:"product_type"`
	ProductImages     []ProductImage `json:"product_images"` //group of pictures
	ProductImage      ProductImage   `json:"product_image"`  //main picture
}

type ProductImage struct {
	gorm.Model
	Image     string `json:"image"`
	ProductID uint   `json:"product_id"`
}

type ProductCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
