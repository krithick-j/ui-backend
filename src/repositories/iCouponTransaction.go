package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveICouponTx(obj models.ICouponTransaction, tx *gorm.DB) error {

	err :=
		tx.
			Create(&obj).
			Error
	return err
}

func GetICouponHistoryByDate(distribId string, fromDate string, toDate string, tx *gorm.DB) ([]models.ICouponTransaction, error) {
	var iCouponHistory []models.ICouponTransaction

	err :=
		tx.
			Model(&models.ICouponTransaction{}).
			Where("distrib_id=? AND created_at BETWEEN ? AND ?", distribId, fromDate, toDate).
			Take(&iCouponHistory).
			Error
	return iCouponHistory, err
}

func GetICouponHistory(distribId string, tx *gorm.DB) ([]models.ICouponTransaction, error) {
	var iCouponHistory []models.ICouponTransaction

	err :=
		tx.
			Model(&models.ICouponTransaction{}).
			Where("distrib_id=?", distribId).
			Limit(20).
			Take(&iCouponHistory).
			Error
	return iCouponHistory, err
}

func GetICouponValueByReferenceNo(referenceNo string, distribId string, tx *gorm.DB) ([]models.ICouponTransaction, error) {
	iCoupons := []models.ICouponTransaction{}
	err :=
		tx.
			Model(&models.ICouponTransaction{}).
			Where("distrib_id=? AND reference=?", distribId, referenceNo).
			Take(&iCoupons).
			Error
	return iCoupons, err
}

func GetICouponsVIDByReference(referenceNo string, distribId string, tx *gorm.DB) ([]string, error) {
	V_IDArr := []string{}
	err :=
		tx.
			Model(&models.ICouponTransaction{}).
			Distinct("v_id").
			Where("distrib_id=? AND reference=?", distribId, referenceNo).
			Pluck("v_id", &V_IDArr).
			Error
	return V_IDArr, err
}
