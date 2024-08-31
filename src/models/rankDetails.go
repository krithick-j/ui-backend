package models

type RankDetails struct {
	ID     uint `gorm:"primarykey"`
	RankId float64
	Type   string
	Target int
}
