package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func AddRspTx(distrib_id string, order_id string, total_rsp float64) *gorm.DB {
	tx := models.RspTransaction{
		DistribId: distrib_id,
		OrderId:   order_id,
		Rsp:       total_rsp,
	}
	result := configs.DB.Create(&tx)
	return result
}

func GetAllRspByDistribId(distrib_id string) ([]int, *gorm.DB) {
	var rsp []int
	result := configs.DB.Model(models.RspTransaction{}).Select("rsp").Where("distrib_id= ?", distrib_id).Find(&rsp)
	return rsp, result
}
func AddDirectBvTx(distrib_id string, order_id string, total_bv float64) *gorm.DB {
	tx := models.RspTransaction{
		DistribId: distrib_id,
		OrderId:   order_id,
		DirectBv:  total_bv,
	}
	result := configs.DB.Create(&tx) //insert into rspTransaction
	return result
}
