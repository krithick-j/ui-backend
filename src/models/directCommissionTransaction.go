package models

import (
	"time"

	"gorm.io/gorm"
)

type DirectCommissionTransaction struct {
	gorm.Model
	DistribId    string    `json:"distrib_id"`
	Value        float64   `json:"value"`
	Reference    string    `json:"reference"`
	RefDistribId string    `json:"ref_distrib_id"` //the comission came from the person
	ActivateDate time.Time `json:"activate_date"`  //21 days after buying product
	ExpiryDate   time.Time `json:"expiry_date"`    //after 6 months
	IsActive     bool      `json:"is_active"`
}
