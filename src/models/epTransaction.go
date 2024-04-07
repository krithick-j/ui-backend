package models

import "gorm.io/gorm"

type EpTransaction struct {
	gorm.Model
	DistribId string `json:"distrib_id"`
	Reference string `json:"reference"`
	Value     int    `json:"value"`
}

