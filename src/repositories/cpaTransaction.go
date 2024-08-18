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
			Model(&models.CpaTransaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("distrib_id=?", distrib_id).
			Take(&value).
			Error
	return value, err
}

func GetAvailableCpaBalance(distrib_id string, tx *gorm.DB) (float64, error) {
	var value float64
	err :=
		tx.
			Model(&models.CpaTransaction{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("distrib_id=? AND is_active=1", distrib_id).
			Take(&value).
			Error
	return value, err
}
