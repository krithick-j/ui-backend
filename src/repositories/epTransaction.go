package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveEpTx(distrib_id string, order_id string, total_ep_value float64) *gorm.DB {

	tx := models.EpTransaction{
		DistribId: distrib_id,
		Reference: order_id,
		Value:     -total_ep_value,
	}

	result := configs.DB.Create(&tx)

	return result
}

func GetEpBalance(distrib_id string) (float64, *gorm.DB) {
	var total float64
	result := configs.DB.Table("ep_transactions").
		Select("sum(value)").
		Where("disrib_id = ?", distrib_id).Scan(&total)
	return total, result
}
