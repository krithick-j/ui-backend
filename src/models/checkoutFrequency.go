package models

import "gorm.io/gorm"

//ChequeFrequencies check for valid cheque draw
type ChequeFrequencies struct {
	gorm.Model
	DistribID string `json:"distrib_id"`
	Frequency int    `json:"frequency"`
	TrackingCenter string    `json:"tracking_center"`
}
