package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveICouponTx(tx models.ICouponTransaction) *gorm.DB {

	result := configs.DB.Create(&tx)
	return result
}

func GetICouponHistory(distribId string) ([]models.ICouponTransaction, error) {
	var iCouponHistory []models.ICouponTransaction
	result := configs.DB.Find(iCouponHistory, "distrib_id=?", distribId)
	return iCouponHistory, result.Error
}
