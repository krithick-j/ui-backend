package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRank struct {
	gorm.Model
	DistribId string     `json:"distrib_id"`
	Year      uint64     `json:"year"`
	Month     time.Month `json:"month"`
	Rank      float64    `json:"rank"`
}
