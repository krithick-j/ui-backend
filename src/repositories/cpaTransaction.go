package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveCpaTransaction(object models.CpaTransaction, tx *gorm.DB) error {

	err :=
		tx.
			Create(&object).
			Error
	return err
}

func GetCpaBalance(distrib_id string, tx *gorm.DB) (float64, error) {
	var value float64
	err :=
		tx.
			Table("cpa_transaction").
			Select("COALESCE(SUM(bv_value), 0)").
			Where("distrib_id", distrib_id).
			Take(&value).
			Error
	return value, err
}
