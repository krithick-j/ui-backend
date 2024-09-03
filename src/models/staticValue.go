package models

import "gorm.io/gorm"

type StaticValue struct {
	gorm.Model
	SocialConn string
}
