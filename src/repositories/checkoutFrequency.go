package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetCheckoutFrequency(distrib_id string) (int, *gorm.DB) {
	var frequency int
	result := configs.DB.Table("cheque_frequencies").Select("frequency").Where("distrib_id", distrib_id).Take(&frequency)
	return frequency, result
}

func IncrementCheckoutFrequency(distrib_id string, place string, frequency int) *gorm.DB {
	result := configs.DB.Where("distrib_id=?", distrib_id).Save(&models.ChequeFrequency{DistribId: distrib_id, Place: place, Frequency: frequency})
	return result
}
