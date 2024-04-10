package repositories

import (
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveBvTransaction(tx models.BvTransaction) *gorm.DB {

	result := configs.DB.Create(&tx)
	return result
}

func GetBvHistoryByTransType(distribId string, transType string, fromDate time.Time, toDate time.Time) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	query := configs.DB.Where("distrib_id=? AND trans_type=?", distribId, transType)
	if !fromDate.IsZero() && !toDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", fromDate, toDate)
	}
	result := query.Find(&BvHistory)
	return BvHistory, result.Error
}

func GetAllBvHistory(distribId string) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	result := configs.DB.Find(BvHistory, "distrib_id=?", distribId)
	return BvHistory, result.Error
}
