package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func UpdateFinalBvToTc(distrib_id string, place string, finalTotalBV int, tx *gorm.DB) error {
	err :=
		tx.
			Model(models.TrackingCenter{}).
			Where("distrib_id=? AND place=?", distrib_id, place).
			Update("bv", finalTotalBV).
			Error
	return err
}

func GetTrackingCenterByDistribId(distrib_id string, trackingCenters []models.TrackingCenter, tx *gorm.DB) ([]models.TrackingCenter, error) {
	err :=
		tx.
			Table("tracking_centers").
			Where("distrib_id", distrib_id).
			Find(&trackingCenters).
			Error
	return trackingCenters, err
}

// func UpdateParentPlaceBv(distrib_id string, placeBv dto.PlaceBv,currentBv float64,  tx *gorm.DB) (float64, error) {
// 	currentBv, _ = GetCurrentBvFromTc(distrib_id, placeBv.Place, currentBv)
// 	updatedBv := currentBv + placeBv.AddBv
// 	err := tx.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("bv", updatedBv)

// 	return err, currentBv
// }

// func UpdateLeftPointPlaceBv(distrib_id string, placeBv dto.PlaceBv) (*gorm.DB, float64) {
// 	var currentLeftPoint float64
// 	currentLeftPoint, _ = GetCurrentLeftPointFromTc(distrib_id, placeBv.Place, currentLeftPoint)
// 	currentLeftPoint += placeBv.AddBv
// 	err := tx.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("left_point", currentLeftPoint)
// 	return err, currentLeftPoint
// }

// func UpdateRightPointPlaceBv(distrib_id string, placeBv dto.PlaceBv) (*gorm.DB, float64) {
// 	var currentRightPoint float64
// 	currentRightPoint, _ = GetCurrentRightPointFromTc(distrib_id, placeBv.Place, currentRightPoint)
// 	currentRightPoint += placeBv.AddBv
// 	err := tx.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("right_point", currentRightPoint)
// 	return err, currentRightPoint
// }

// func UpdateTrackingCenter(leftPoint int, rightPoint int, distrib_id string, place string) *gorm.DB {

// 	err := tx.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, place).Updates(map[string]interface{}{"left_point": leftPoint, "right_point": rightPoint})
// 	return err
// }

// func GetCurrentBvFromTc(distrib_id string, place string, bv float64) (float64, *gorm.DB) {
// 	err := tx.Table("tracking_centers").Select("bv").Where("distrib_id=? AND place=?", distrib_id, place).Take(&bv)
// 	return bv, err
// }

// func GetCurrentLeftPointFromTc(distrib_id string, place string, leftPoint float64) (float64, *gorm.DB) {
// 	err := tx.Table("tracking_centers").Select("left_point").Where("distrib_id=? AND place=?", distrib_id, place).Take(&leftPoint)
// 	return leftPoint, err
// }

// func GetCurrentRightPointFromTc(distrib_id string, place string, rightPoint float64) (float64, *gorm.DB) {
// 	err := tx.Table("tracking_centers").Select("right_point").Where("distrib_id=? AND place=?", distrib_id, place).Take(&rightPoint)
// 	return rightPoint, err
// }

func GetTrackingCenter(distrib_id, place string, tx *gorm.DB) (*models.TrackingCenter, error) {
	tc := new(models.TrackingCenter)
	err :=
		tx.
			First(tc, "distrib_id = ? AND place = ?", distrib_id, place).
			Error
	return tc, err
}

func GetAllTrackingCenters(distrib_id string, tx *gorm.DB) ([]models.TrackingCenter, error) {
	var tcs []models.TrackingCenter
	err := tx.Find(&tcs, "distrib_id = ?", distrib_id).Error
	return tcs, err
}

func ActivateTC(distrib_id string, place string, tx *gorm.DB) error {
	err :=
		tx.
			Table("tracking_centers").
			Where("distrib_id=? AND place=?", distrib_id, place).
			Updates(map[string]interface{}{"is_active": 1}).
			Error
	return err
}
