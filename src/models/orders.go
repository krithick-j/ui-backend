package models

import (
	"gorm.io/gorm"
)

type OrdersHeader struct {
	gorm.Model
	DistribId     string
	OrderId       string
	SubTotal      float64
	TotalSandH    float64
	TotalAmount   float64
	TotalQuantity float64
	TotalBV       int
	ContactName   string
	ContactEmail  string
	Address       string
	City          string
	District      string
	State         string
	ZipCode       uint64
	Country       string
	HomePhoneNo   string
	MobilePhoneNo string
	OrdersLiner   []OrdersLiner
}

type OrdersLiner struct {
	gorm.Model
	OrdersHeaderID uint
	Name           string
	Quantity       uint
	UnitPrice      uint64
	ProductType    string
	TypeValue      float64
	SubTotal       float64
	SandH          float64
}
