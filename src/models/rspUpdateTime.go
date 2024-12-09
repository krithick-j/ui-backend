package models

import (
	"time"

	"gorm.io/gorm"
)

type RspUpdateTime struct {
	gorm.Model
	RspUpdate time.Time
}
