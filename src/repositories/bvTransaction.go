package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func SaveBvTransaction(tx models.BvTransaction) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
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
