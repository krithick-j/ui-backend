package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetRankValueByDistribId(distrib_id string) (float64, *gorm.DB) {
	var rank float64
	result := configs.DB.Model(&models.UserRank{}).Select("rank").Where("distrib_id=?", distrib_id).Take(&rank)
	return rank, result
}

func SaveRankValue(RankObject models.UserRank) error {
	err := configs.DB.Model(&models.UserRank{}).Create(&RankObject).Error
	return err 
}