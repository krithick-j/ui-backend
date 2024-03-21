package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetCheckoutFrequency(distrib_id string) (int, *gorm.DB) {
	var frequency int
	result := configs.DB.Table("checkout_frequencies").Select("frequency").Where("distrib_id", distrib_id).Take(&frequency)
	return frequency, result
}

func IncrementCheckoutFrequency(distrib_id string, frequency int) *gorm.DB {
	result := configs.DB.Model(models.CheckoutFrequency{}).Where("distrib_id=?", distrib_id).Update("frequency", frequency+1)
	return result
}
