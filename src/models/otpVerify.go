package models

import (
	"gorm.io/gorm"
)

type OTPVerify struct {
	gorm.Model
	Type   string `json:"type"` // email | mobile
	Value  string `json:"value"`
	OTP    string `json:"OTP"`
	Status string `json:"status" gorm:"default:unverified"` // unverified | verified
}
