package repositories

import (
	"time"

	"gorm.io/gorm"
)

// Get the last Updated Cron tab from rank_redeem.py
func GetRspTime(tx *gorm.DB) (time.Time, error) {
	var rspUpdateTime time.Time
	err :=
		tx.
			Table("rsp_update_times").
			Select("rsp_update").
			Take(&rspUpdateTime).
			Error
	return rspUpdateTime, err
}
