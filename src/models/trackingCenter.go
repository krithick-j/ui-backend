package models

import "gorm.io/gorm"

type TrackingCenter struct {
	gorm.Model
	DistribID      string `json:"distrib_id"`
	Name           string `json:"name"`
	PDistribId     string `json:"p_distrib_id"`
	PPlace         string `json:"p_place"`
	LeftDistribID  string `json:"left_distrib_id"`
	LeftPlace      string `json:"left_place"`
	RightDistribID string `json:"right_distrib_id"`
	RightPlace     string `json:"right_place"`
	Place          string `json:"place"`
	Bv             int    `json:"bv"`
	LeftPoint      int    `json:"left_point" gorm:"default:0"`
	RightPoint     int    `json:"right_point" gorm:"default:0"`
	IsActive       bool   `json:"is_active" gorm:"default:0"`
}

type TCBv struct {
	DistribID string 
	Side      string
	BValue    int
}
