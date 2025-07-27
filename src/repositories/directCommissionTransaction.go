package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveDirectCommissionTransaction(obj models.DirectCommissionTransaction, tx *gorm.DB) error {
	err :=
		tx.
			Create(&obj).
			Error
	return err
}

func GetDirectCommissionActiveValueByDistribId(ref_distrib_id string, tx *gorm.DB) (float64, error) {
	var value float64
	err :=
		tx.
			Model(&models.DirectCommissionTransaction{}).
			Select("COALESCE(sum(value), 0)").
			Where("ref_distrib_id =?  AND is_active = 1", ref_distrib_id).
			Take(&value).
			Error
	return value, err
}

func GetDirectCommissionAllValueByDistribId(ref_distrib_id string, tx *gorm.DB) (float64, error) {
    var value float64
    err :=
        tx.
            Model(&models.DirectCommissionTransaction{}).
            Select("COALESCE(sum(value), 0)").
            Where("ref_distrib_id = ? AND expiry_date > NOW()", ref_distrib_id).
            Take(&value).
            Error
    return value, err
}