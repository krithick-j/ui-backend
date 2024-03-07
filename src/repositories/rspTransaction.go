package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func AddRspTx(distrib_id string, order_id string, total_rsp int) {
	tx := models.RspTransaction{
		DistribId: distrib_id,
		OrderId:   order_id,
		Rsp:       total_rsp,
	}
	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}

func GetAllRspByDistribId(distrib_id string) ([]int, *gorm.DB) {
	var rsp []int
	result := configs.DB.Model(models.RspTransaction{}).Select("rsp").Where("distrib_id= ?", distrib_id).Find(&rsp)
	return rsp, result
}
func AddDirectBvTx83(distrib_id string, order_id string, total_bv int) {
	tx := models.RspTransaction{
		DistribId: distrib_id,
		OrderId:   order_id,
		DirectBv:  total_bv,
	}
	result := configs.DB.Create(&tx) //insert into rspTransaction
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}