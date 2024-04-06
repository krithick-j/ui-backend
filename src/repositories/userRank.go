package repositories

import (
	"ui-back-end/configs"

	"gorm.io/gorm"
)

func GetRankValueByDistribId(distrib_id string) (float32, *gorm.DB) {
	var rank float32
	result := configs.DB.Table("user_rank").Select("rank").Where("distrib_id=?", distrib_id).Take(&rank)
	return rank, result
}
