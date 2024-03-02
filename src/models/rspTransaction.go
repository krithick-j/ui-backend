package models

import (
	"time"

	"gorm.io/gorm"
)

type RspTransaction struct {
	gorm.Model
	DistribId string    `json:"distrib_id"`
	OrderId   string    `json:"order_id"`
	Rsp       int       `json:"rsp"`
	Date      time.Time `json:"date"`
}
