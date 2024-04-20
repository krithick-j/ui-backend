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

func GetICouponHistoryByDate(distribId string, fromDate string, toDate string) ([]models.ICouponTransaction, error) {
	var iCouponHistory []models.ICouponTransaction

	result := configs.DB.Model(&models.ICouponTransaction{}).Where("distrib_id=? AND created_at BETWEEN ? AND ?", distribId, fromDate, toDate).Take(&iCouponHistory)
	return iCouponHistory, result.Error
}

func GetICouponHistory(distribId string) ([]models.ICouponTransaction, error) {
	var iCouponHistory []models.ICouponTransaction

	result := configs.DB.Model(&models.ICouponTransaction{}).Where("distrib_id=?", distribId).Limit(20).Take(&iCouponHistory)
	return iCouponHistory, result.Error
}
