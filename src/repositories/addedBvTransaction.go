package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
)

func GetAddedBvByOrderIdAndDistribId(orderId string, distribId string) ([]dto.PlaceBv, error) {
	placeBv := []dto.PlaceBv{}
	err := configs.DB.Model(&models.AddedBv{}).Select("place as Place, value as AddBv").Where("distrib_id = ? AND reference_no =?", distribId, orderId).Scan(&placeBv).Error
	return placeBv, err
}

// Saves the object
func SaveAddedbv(obj *models.AddedBv) error {
	err := configs.DB.Create(&obj).Error
	return err
}
