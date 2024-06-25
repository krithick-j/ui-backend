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

func GetICouponValueByReferenceNo(referenceNo string, distribId string) ([]models.ICouponTransaction, error) {
	iCoupons := []models.ICouponTransaction{}
	err := configs.DB.Model(&models.ICouponTransaction{}).Where("distrib_id=? AND reference=?", distribId, referenceNo).Take(&iCoupons).Error
	return iCoupons, err
}

func GetICouponsVIDByReference(referenceNo string, distribId string) ([]string, error) {
	V_IDArr := []string{}
	err := configs.DB.Model(&models.ICouponTransaction{}).Distinct("v_id").Where("distrib_id=? AND reference=?", distribId, referenceNo).Pluck("v_id", &V_IDArr).Error
	return V_IDArr, err
}
