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

func GetDirectCommissionValueByDistribId(distribId string) (float64, *gorm.DB) {
	var value float64
	result := configs.DB.Model(&models.DirectCommissionTransaction{}).Select("sum(value)").Where("distrib_id=?", distribId).Take(&value)
	return value, result
}
