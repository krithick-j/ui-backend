package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func UpdateFinalBvToTc(distrib_id string, place string, finalTotalBV int) *gorm.DB {
	result := configs.DB.Model(models.TrackingCenter{}).Where("distrib_id=? AND place=?", distrib_id, place).Update("bv", finalTotalBV)
	return result
}
