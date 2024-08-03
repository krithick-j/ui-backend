package models

import (
	"gorm.io/gorm"
)

type RspTransaction struct {
	gorm.Model
	DistribId string  `json:"distrib_id"`
	OrderId   string  `json:"order_id"`
	Rsp       float64 `json:"rsp" gorm:"default:0"`
	DirectBv  float64 `json:"direct_bv" gorm:"default:0"`
}
