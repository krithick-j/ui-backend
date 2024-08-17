package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveEpTx(distrib_id string, order_id string, total_ep_value float64, tx *gorm.DB) error {

	obj := models.EpTransaction{
		DistribId: distrib_id,
		Reference: order_id,
		Value:     -total_ep_value,
	}

	err := tx.Create(&obj).Error

	return err
}

func GetEpBalance(distrib_id string, tx *gorm.DB) (float64, error) {
	var total float64
	err := tx.Table("ep_transactions").
		Select("COALESCE(sum(value), 0)").
		Where("distrib_id = ?", distrib_id).
		Scan(&total).
		Error
	return total, err
}
