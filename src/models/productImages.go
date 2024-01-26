package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	Image     string
	ProductID uint
}
