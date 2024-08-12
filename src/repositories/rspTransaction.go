package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func AddRsp(distrib_id string, order_id string, total_rsp float64) error {
	obj := models.RspTransaction{
		DistribId: distrib_id,
		OrderId:   order_id,
		Rsp:       total_rsp,
	}
	err := configs.DB.Create(&obj).Error
	return err
}

// Total Personal RSP sum of distributor id
func GetPersonalRspSumByDistribId(distrib_id string, tx *gorm.DB) (float64, error) {
	var totalRsp float64
	err := tx.Model(models.RspTransaction{}).Select("COALESCE(SUM(rsp),0)").Where("distrib_id= ?", distrib_id).Find(&totalRsp).Error
	return totalRsp, err
}
func AddDirectBvTx(referral_distrib_id string, order_id string, total_bv float64) *gorm.DB {
	tx := models.RspTransaction{
		DistribId: referral_distrib_id,
		OrderId:   order_id,
		DirectBv:  total_bv,
	}
	result := configs.DB.Create(&tx) //insert into rspTransaction
	return result
}

func GetDirectBvByDistribID(distrib_id string, tx *gorm.DB) (float64, error) {
	var totalRsp float64
	err := tx.Model(models.RspTransaction{}).Select("SUM(direct_bv)").Where("distrib_id= ?", distrib_id).Find(&totalRsp).Error
	return totalRsp, err
}
