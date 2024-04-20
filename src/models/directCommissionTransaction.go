package models

import "gorm.io/gorm"

type DirectCommissionTransaction struct {
	gorm.Model
	DistribId string  `json:"distrib_id"`
	Value     float64 `json:"value"`
	Reference string  `json:"reference"`
	FromDistribId string `json:"from_distrib_id"` //the comission came from the person
}
