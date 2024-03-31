package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func UpdateFinalBvToTc(distrib_id string, place string, finalTotalBV int) *gorm.DB {

	result := configs.DB.Model(models.TrackingCenter{}).Where("distrib_id=? AND place=?", distrib_id, place).Update("bv", finalTotalBV)
	return result
}

func GetTrackingCenterByDistribId(distrib_id string, trackingCenters []models.TrackingCenter) ([]models.TrackingCenter, *gorm.DB) {
	result := configs.DB.Table("tracking_centers").Where("distrib_id", distrib_id).Find(&trackingCenters)
	return trackingCenters, result
}

func UpdateParentPlaceBv(distrib_id string, placeBv dto.PlaceBv) (*gorm.DB, int) {
	var currentBv int
	fmt.Println("place", placeBv.Place)
	currentBv, _ = GetCurrentBvFromTc(distrib_id, placeBv.Place, currentBv)
	updatedBv := currentBv + placeBv.AddBv
	fmt.Println("Hello this is current bv ", currentBv)
	fmt.Println("place bv", placeBv)
	result := configs.DB.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("bv", updatedBv)
	if result.Error != nil {
		fmt.Println("ERROR IN UPDATING BV", result.Error)
	}
	return result, currentBv
}

func UpdateLeftPointPlaceBv(distrib_id string, placeBv dto.PlaceBv) (*gorm.DB, int) {
	var currentLeftPoint int
	currentLeftPoint, _ = GetCurrentLeftPointFromTc(distrib_id, placeBv.Place, currentLeftPoint)
	currentLeftPoint += placeBv.AddBv
	result := configs.DB.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("left_point", currentLeftPoint)
	return result, currentLeftPoint
}

func UpdateRightPointPlaceBv(distrib_id string, placeBv dto.PlaceBv) (*gorm.DB, int) {
	var currentRightPoint int
	currentRightPoint, _ = GetCurrentRightPointFromTc(distrib_id, placeBv.Place, currentRightPoint)
	currentRightPoint += placeBv.AddBv
	result := configs.DB.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, placeBv.Place).Update("right_point", currentRightPoint)
	return result, currentRightPoint
}

func UpdateTrackingCenter(leftPoint int, rightPoint int, distrib_id string, place string) *gorm.DB {

	result := configs.DB.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, place).Updates(map[string]interface{}{"left_point": leftPoint, "right_point": rightPoint})
	return result
}

func GetCurrentBvFromTc(distrib_id string, place string, bv int) (int, *gorm.DB) {
	result := configs.DB.Table("tracking_centers").Select("bv").Where("distrib_id=? AND place=?", distrib_id, place).Take(&bv)
	return bv, result
}

func GetCurrentLeftPointFromTc(distrib_id string, place string, leftPoint int) (int, *gorm.DB) {
	result := configs.DB.Table("tracking_centers").Select("left_point").Where("distrib_id=? AND place=?", distrib_id, place).Take(&leftPoint)
	return leftPoint, result
}

func GetCurrentRightPointFromTc(distrib_id string, place string, rightPoint int) (int, *gorm.DB) {
	result := configs.DB.Table("tracking_centers").Select("right_point").Where("distrib_id=? AND place=?", distrib_id, place).Take(&rightPoint)
	return rightPoint, result
}

func GetTrackingCenter(distrib_id, place string) *models.TrackingCenter {
	tc := new(models.TrackingCenter)
	configs.DB.First(tc, "distrib_id = ? AND place = ?", distrib_id, place)
	return tc
}
