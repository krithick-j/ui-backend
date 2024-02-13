package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	DistribId string
	ProductID uint
}
