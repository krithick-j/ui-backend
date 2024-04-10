package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveBvTransaction(tx models.BvTransaction) *gorm.DB {

	result := configs.DB.Create(&tx)
	return result
}

func GetBvHistory(distribId string, transType string) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	result := configs.DB.Find(BvHistory, "distrib_id=? AND trans_type=?", distribId, transType)
	return BvHistory, result.Error
}

func GetAllBvHistory(distribId string) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	result := configs.DB.Find(BvHistory, "distrib_id=?", distribId)
	return BvHistory, result.Error
}
