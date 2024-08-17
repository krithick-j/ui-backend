package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetCheckoutFrequency(distrib_id string, tx *gorm.DB) (int, error) {
	var frequency int
	err :=
		tx.
			Table("cheque_frequencies").
			Select("frequency").
			Where("distrib_id", distrib_id).
			Take(&frequency).
			Error
	return frequency, err
}

func IncrementCheckoutFrequency(distrib_id string, place string, frequency int, tx *gorm.DB) error {
	err :=
		tx.
			Model(models.ChequeFrequency{}).
			Where("distrib_id=?", distrib_id).
			Update("frequency", frequency+1).
			Error
	return err
}

func CreateCheckoutFrequency(distrib_id string, place string, tx *gorm.DB) error {
	obj := models.ChequeFrequency{
		DistribId: distrib_id,
		Frequency: 0,
	}

	err := tx.Create(&obj).Error
	return err
}
