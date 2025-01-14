package repositories

import (
	"fmt"
	"time"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveBvTransaction(object models.BvTransaction, tx *gorm.DB) error {

	err :=
		tx.
			Create(&object).
			Error
	return err
}

func GetBvHistoryByTransType(distribId string, transType string, fromDate time.Time, toDate time.Time, tx *gorm.DB) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	query :=
		tx.
			Where("distrib_id=? AND trans_type=?", distribId, transType)
	if !fromDate.IsZero() && !toDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", fromDate, toDate)
	}
	err := query.Find(&BvHistory).Limit(5).Error
	return BvHistory, err
}

func GetAllBvHistory(distribId string, tx *gorm.DB) ([]models.BvTransaction, error) {
	var BvHistory []models.BvTransaction
	err :=
		tx.
			Find(BvHistory, "distrib_id=?", distribId).
			Error
	return BvHistory, err
}

func GetBVforTC(distrib_id string, tc string, tx *gorm.DB) ([]models.TCBv, error) {
	tcbv := []models.TCBv{}
	err :=
		tx.
			Table("bv_transactions").
			Select("side, sum(bv_value) as BValue").
			Where("distrib_id = ? AND place = ? AND is_active = 1 ", distrib_id, tc).
			Group("side").
			Scan(&tcbv).
			Error
	return tcbv, err
}

// is_active both 0 and 1 will be added for bv counter
func GetBVforTCByDate(distrib_id string, tc string, fromDate time.Time, toDate time.Time, tx *gorm.DB) ([]models.TCBv, error) {
	tcbv := []models.TCBv{}
	err :=
		tx.
			Table("bv_transactions").
			Select("side, sum(bv_value) as BValue").
			Where("distrib_id = ? AND place = ?", distrib_id, tc).
			Where("date BETWEEN ? AND ?", fromDate, toDate).
			Group("side").
			Scan(&tcbv).
			Error
	return tcbv, err
}

func GetBVforTCOneRow(distrib_id string, place string, tx *gorm.DB) (models.TCBvOneRow, error) {
	tcbv := models.TCBvOneRow{}
	err :=
		tx.
			Table("bv_transactions").
			Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
			Where("distrib_id = ? AND place = ? AND is_active = 1 ", distrib_id, place).
			Take(&tcbv).
			Error
	return tcbv, err
}

// Get Tracking Center BV by Date
//Used only in Bv Counter
func GetBVforTCOneRowByDate(distrib_id string, tc string, fromDate time.Time, toDate time.Time, tx *gorm.DB) (models.TCBvOneRow, error) {
	tcbv := models.TCBvOneRow{}
	err :=
		tx.
			Table("bv_transactions").
			Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
			Where("distrib_id = ? AND place = ?", distrib_id, tc).
			Where("date BETWEEN ? AND ?", fromDate, toDate).
			Where("trans_type='product'").
			Take(&tcbv).
			Error
	return tcbv, err
}

// Get Tracking Center BV by Date Generic function
// set parameter fromDate and toDate as empty string to remove activate_date filter
// set toDate as empty string to to get date from till last
func GetBVforTCOneRowByDateGeneric(distrib_id string, tc string, is_active bool, fromDate string, toDate string, tx *gorm.DB) (models.TCBvOneRow, error) {
	tcbv := models.TCBvOneRow{}
	query :=
		tx.
			Table("bv_transactions").
			Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
			Where("distrib_id = ? AND place = ? AND is_active = ?", distrib_id, tc, is_active)

	if fromDate != "" && toDate != "" {
		query =
			query.
				Where("activate_date BETWEEN ? AND ?", fromDate, toDate)
	} else if toDate != "" {
		query =
			query.
				Where("activate_date >= ?", fromDate)
	}
	err :=
		query.
			Take(&tcbv).
			Error
	return tcbv, err
}

func GetBVforTCOneRowByDateForCheque(distrib_id string, tc string, fromDate string, toDate string, tx *gorm.DB) (models.TCBvOneRow, error) {
	tcbv := models.TCBvOneRow{}
	query :=
		tx.
			Table("bv_transactions").
			Select("SUM(IF(side='bv', bv_value, 0)) as b_value, SUM(IF(side='left',bv_value, 0)) as l_value, SUM(IF(side='right', bv_value, 0)) as r_value").
			Where("distrib_id = ? AND place = ? AND ((is_active = 1 AND trans_type='product') OR (trans_type='cheque'))", distrib_id, tc)

	if fromDate != "" && toDate != "" {
		query =
			query.
				Where("activate_date BETWEEN ? AND ?", fromDate, toDate)
	} else if toDate != "" {
		query =
			query.
				Where("activate_date >= ?", fromDate)
	}
	err :=
		query.
			Take(&tcbv).
			Error
	return tcbv, err
}

// Get Sum of Bv transactions of an individual distributor with trans type product
func GetBvSumByDistribId(distribId, transType string, tx *gorm.DB) (int, error) {
	var bvSum int

	err :=
		tx.
			Table("bv_transactions").
			Select("COALESCE(SUM(bv_value), 0)").
			Where("distrib_id=? AND trans_type=?", distribId, transType).
			Scan(&bvSum).
			Error
	return bvSum, err
}

// This function will give you the count of cheque taken by the distrib id in a Month
func GetStepByDistribId(distribId string, tx *gorm.DB) (int64, error) {
	var count int64

	firstOfMonth := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	fmt.Print("month date", firstOfMonth)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	fmt.Print("month date last:", lastOfMonth)
	query :=
		tx.
			Table("bv_transactions").
			Where("distrib_id=? AND trans_type='cheque'", distribId)
	query = query.Where("created_at BETWEEN ? AND ?", firstOfMonth, lastOfMonth)
	err := query.Count(&count).Error
	return count, err
}
