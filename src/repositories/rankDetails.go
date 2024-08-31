package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetRankDetailsByDistribId(rankId float64, tx *gorm.DB) ([]models.RankDetails, error) {
	var RankDetails []models.RankDetails
	err := tx.Table("rank_details").Where("rank_id=?", rankId).Find(&RankDetails).Error
	return RankDetails, err
}
