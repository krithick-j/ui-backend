package models

import (
	"time"

	"gorm.io/gorm"
)

type RspTransaction struct {
	gorm.Model
	DistribId string    `json:"distrib_id"`
	OrderId   string    `json:"order_id"`
	Rsp       int       `json:"rsp gorm:"default:0"`
	Date      time.Time `json:"date"`
	DirectBv  int       `json:"direct_bv" gorm:"default:0"`
}
