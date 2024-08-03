package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTotalRspByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

	rspTotal, err := repositories.GetPersonalRspSumByDistribId(DistribID, tx)
	if err != nil {
		configs.Log.Errorln("Error on calling GetPersonalRspSumByDistribId from GetTotalRspByDistribId service")
		return fiber.Map{"data": rspTotal}, fiber.StatusInternalServerError
	}
	configs.Log.Infoln("GetTotalRspByDistribId service function complete")

	return fiber.Map{"data": rspTotal}, fiber.StatusOK
}

func GetGroupRspByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

	//group rsp pananum
	var rspSum float64
	// first avan referral distrib id eduthu avanoda rsp edukanum
	referredDistribIds, err := repositories.GetReferredUsersByDistribId(DistribID, tx)
	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": rspSum}, fiber.StatusOK
	}

	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId from test service fn", err.Error())
	}

	for _, referredDistribId := range referredDistribIds {
		fmt.Println("referred distrib id", DistribID)
		totalRspForOneDistrib, err := repositories.GetPersonalRspSumByDistribId(referredDistribId, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetPersonalRspSumByDistribIdt from test service fn", err.Error())
		}

		rspSum += totalRspForOneDistrib
		//suppose IN-00001 is distrib id and referrerid is also IN-00001 then, it will loop continuously right ? So I am continuing to next distrib id
		if DistribID == referredDistribId {
			continue
		}

		recursiveRspSum, status := Test(referredDistribId)
		if status != fiber.StatusOK {
			tx.Rollback()
			return recursiveRspSum, status
		}
		rspSum += recursiveRspSum["data"].(float64)
	}

	return fiber.Map{"data": rspSum}, fiber.StatusOK
}
