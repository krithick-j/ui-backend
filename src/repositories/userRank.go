package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetRankValueByDistribId(distrib_id string, tx *gorm.DB) (float64, error) {
	var rank float64
	err := tx.Model(&models.UserRank{}).Select("rank").Where("distrib_id=?", distrib_id).Take(&rank).Error
	return rank, err
}
