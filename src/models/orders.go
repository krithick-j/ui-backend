package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	VId []string `json:"v_id"`
	
}