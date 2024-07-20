package repositories

import (
	"fmt"
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
	result := query.Find(&BvHistory).Limit(5)
	return BvHistory, result.Error
}

func GetAllBvHistory(distribId string) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	result := configs.DB.Find(BvHistory, "distrib_id=?", distribId)
	return BvHistory, result.Error
}

func GetBVforTC(distrib_id string, tc string) ([]models.TCBv, *gorm.DB) {
	tcbv := []models.TCBv{}
	result := configs.DB.Table("bv_transactions").
		Select("side, sum(bv_value) as BValue").
		Where("distrib_id = ? AND place = ? AND is_active = 1 ", distrib_id, tc).
		Group("side").
		Scan(&tcbv)
	return tcbv, result
}

// is_active both 0 and 1 will be added for bv counter
func GetBVforTCByDate(distrib_id string, tc string, fromDate time.Time, toDate time.Time) ([]models.TCBv, *gorm.DB) {
	tcbv := []models.TCBv{}
	print("hi from tc")
	result := configs.DB.Table("bv_transactions").
		Select("side, sum(bv_value) as BValue").
		Where("distrib_id = ? AND place = ?", distrib_id, tc).
		Where("date BETWEEN ? AND ?", fromDate, toDate).
		Group("side").
		Scan(&tcbv)
	return tcbv, result
}

func GetBVforTCOneRow(distrib_id string, place string) (models.TCBvOneRow, *gorm.DB) {
	tcbv := models.TCBvOneRow{}
	result := configs.DB.Table("bv_transactions").
		Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
		Where("distrib_id = ? AND place = ? AND is_active = 1 ", distrib_id, place).
		Take(&tcbv)
	return tcbv, result
}

// Get Tracking Center BV by Date
func GetBVforTCOneRowByDate(distrib_id string, tc string, fromDate time.Time, toDate time.Time) (models.TCBvOneRow, *gorm.DB) {
	tcbv := models.TCBvOneRow{}
	fmt.Print("hi from row")
	result := configs.DB.Table("bv_transactions").
		Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
		Where("distrib_id = ? AND place = ? AND is_active = 1 ", distrib_id, tc).
		Where("activate_date BETWEEN ? AND ?", fromDate, toDate).
		Take(&tcbv)
	return tcbv, result
}

// Get Tracking Center BV by Date Generic function
// set parameter fromDate and toDate as empty string to remove activate_date filter
// set toDate as empty string to to get date from till last
func GetBVforTCOneRowByDateGeneric(distrib_id string, tc string, is_active bool, fromDate string, toDate string) (models.TCBvOneRow, *gorm.DB) {
	tcbv := models.TCBvOneRow{}
	fmt.Print("hi from row")
	query := configs.DB.Table("bv_transactions").
		Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
		Where("distrib_id = ? AND place = ? AND is_active = ?", distrib_id, tc, is_active)

	if fromDate != "" && toDate != "" {
		query = query.Where("activate_date BETWEEN ? AND ?", fromDate, toDate)
	} else if toDate != "" {
		query = query.Where("activate_date >= ?", fromDate)
	}
	result := query.Take(&tcbv)
	return tcbv, result
}

// Get Sum of Bv transactions of an individual distributor with trans type product
func GetBvSumByDistribId(distribId, transType string) (int, error) {
	var bvSum int

	err := configs.DB.Table("bv_transactions").Select("COALESCE(SUM(bv_value), 0)").Where("distrib_id=? AND trans_type=?", distribId, transType).Scan(&bvSum).Error
	return bvSum, err
}
