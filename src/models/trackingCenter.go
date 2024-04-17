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
	IsActive       bool   `json:"is_active" gorm:"default:0"`
}

type TCBv struct {
	Side   string `json:"side"`
	BValue int    `json:"b_value"`
}

type TCBvOneRow struct {
	BValue int `json:"b_value"`
	LValue int `json:"l_value"`
	RValue int `json:"r_value"`
}
