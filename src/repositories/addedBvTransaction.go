package repositories

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetAddedBvByOrderIdAndDistribId(orderId string, distribId string, tx *gorm.DB) ([]dto.PlaceBv, error) {
	placeBv := []dto.PlaceBv{}
	err :=
		tx.
			Model(&models.AddedBv{}).
			Select("place as Place, value as AddBv").
			Where("distrib_id = ? AND reference_no =?", distribId, orderId).
			Scan(&placeBv).
			Error
	return placeBv, err
}

// Saves the object
func SaveAddedbv(obj *models.AddedBv, tx *gorm.DB) error {
	err := tx.Create(&obj).Error
	return err
}
