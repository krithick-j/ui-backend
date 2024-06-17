package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveDirectCommissionTransaction(obj models.DirectCommissionTransaction) *gorm.DB {
	result := configs.DB.Create(&obj)
	return result
}

func GetDirectCommissionActiveValueByDistribId(ref_distrib_id string) (float64, *gorm.DB) {
	var value float64
	result := configs.DB.Model(&models.DirectCommissionTransaction{}).Select("COALESCE(sum(value), 0)").Where("ref_distrib_id =?  AND is_active = 1", ref_distrib_id).Take(&value)
	return value, result
}

func GetDirectCommissionAllValueByDistribId(ref_distrib_id string) (float64, *gorm.DB) {
	var value float64
	result := configs.DB.Model(&models.DirectCommissionTransaction{}).Select("COALESCE(sum(value), 0)").Where("ref_distrib_id =?", ref_distrib_id).Take(&value)
	return value, result
}
