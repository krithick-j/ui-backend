package models

import "gorm.io/gorm"

type EnquiryType struct {
	gorm.Model
	Name string `json:"name"`
	AdminName string `json:"admin_name"`
}
