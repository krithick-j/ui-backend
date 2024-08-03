package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// func Test(distribId string) (fiber.Map, int) {

// 	configs.Log.Infof("Test service completed")
// 	return fiber.Map{"data": ""}, 200
// }

func Test(distribId string) (fiber.Map, int) {
	tx := configs.DB.Begin()
	//group rsp pananum
	var rspSum float64
	// first avan referral distrib id eduthu avanoda rsp edukanum
	referredDistribIds, err := repositories.GetReferredUsersByDistribId(distribId, tx)
	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": rspSum}, fiber.StatusOK
	}

	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId from test service fn", err.Error())
	}

	for _, referredDistribId := range referredDistribIds {
		fmt.Println("referred distrib id", distribId)
		totalRspForOneDistrib, err := repositories.GetPersonalRspSumByDistribId(referredDistribId, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetPersonalRspSumByDistribIdt from test service fn", err.Error())
		}

		rspSum += totalRspForOneDistrib
		//suppose IN-00001 is distrib id and referrerid is also IN-00001 then, it will loop continuously right ? So I am continuing to next distrib id
		if distribId == referredDistribId {
			continue
		}

		recursiveRspSum, status := Test(referredDistribId)
		if status != fiber.StatusOK {
			tx.Rollback()
			return recursiveRspSum, status
		}
		rspSum += recursiveRspSum["data"].(float64)
	}

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error committing transaction:", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": rspSum}, fiber.StatusOK
}
