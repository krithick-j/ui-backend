package models

import (
	"gorm.io/gorm"
)

type OrdersHeader struct {
	gorm.Model
	DistribId          string
	OrderId            string
	DeliveryStatus     string
	CourierName        string
	ShipmentTrackingNo string
	DeliveredAt        string
	SubTotal           float64
	TotalSandH         float64
	TotalAmount        float64
	TotalQuantity      float64
	TotalTypeValue     float64
	ProductType        string
	ContactName        string
	ContactEmail       string
	Address            string
	City               string
	District           string
	State              string
	ZipCode            uint64
	Country            string
	HomePhoneNo        string
	MobilePhoneNo      string
	OrdersLiner        []OrdersLiner
}

type OrdersLiner struct {
	gorm.Model
	OrdersHeaderID uint
	ProductID      uint
	ProductImage   string
	Name           string
	Quantity       uint
	UnitPrice      float64
	ProductType    string
	TypeValue      float64
	SubTotal       float64
	SandH          float64
}
