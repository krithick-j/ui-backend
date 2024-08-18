package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPersonalRspByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

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
	//avan rsp mattum add aaga koodathu
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId from test service fn", err.Error())
	}
	fmt.Println("referred distrib ids", referredDistribIds)

	for _, referredDistribId := range referredDistribIds {
		//suppose IN-00001 is distrib id and referrerid is also IN-00001 then, it will loop continuously right ? So I am continuing to next distrib id

		if DistribID == referredDistribId {
			continue
		}
		fmt.Println("referred distrib id is--->", referredDistribId, "distrib id is --->", DistribID)

		totalRspForOneDistrib, err := repositories.GetPersonalRspSumByDistribId(referredDistribId, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetPersonalRspSumByDistribIdt from test service fn", err.Error())
		}

		rspSum += totalRspForOneDistrib

		recursiveRspSum, status := GetGroupRspByDistribId(referredDistribId, tx)
		if status != fiber.StatusOK {
			tx.Rollback()
			return recursiveRspSum, status
		}
		rspSum += recursiveRspSum["data"].(float64)
	}

	return fiber.Map{"data": rspSum}, fiber.StatusOK
}

func GetDirectBvByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

	directBv, err := repositories.GetDirectBvByDistribID(DistribID, tx)
	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": directBv}, fiber.StatusOK
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetDirectBvByDistribID from test service fn", err.Error())
	}

	return fiber.Map{"data": directBv}, fiber.StatusOK
}

func GetTotalStepByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

	step, err := repositories.GetStepByDistribId(DistribID, tx)
	//Send 0 if record not found
	if err == gorm.ErrRecordNotFound {
		return utils.SuccessMessage(step, fiber.StatusOK)
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetStepByDistribId", "GetTotalStepByDistribId", fiber.StatusInternalServerError, tx)
	}
	step = step / 2
	return utils.SuccessMessage(step, fiber.StatusOK)
}

func GetRspValuesByDistribID(DistribID string, tx *gorm.DB) (fiber.Map, int) {
	directBvResult, status := GetDirectBvByDistribId(DistribID, tx)
	if status != fiber.StatusOK {
		return directBvResult, status
	}
	directBv := directBvResult["data"]

	personalRspResult, status := GetPersonalRspByDistribId(DistribID, tx)
	if status != fiber.StatusOK {
		return personalRspResult, status
	}
	personalRsp := personalRspResult["data"]

	groupRspResult, status := GetGroupRspByDistribId(DistribID, tx)
	if status != fiber.StatusOK {
		return groupRspResult, status
	}
	groupRsp := groupRspResult["data"]

	stepResult, status := GetGroupRspByDistribId(DistribID, tx)
	if status != fiber.StatusOK {
		return stepResult, status
	}
	step := stepResult["data"]

	response := fiber.Map{
		"direct_bv":    directBv,
		"personal_rsp": personalRsp,
		"group_rsp":    groupRsp,
		"step":         step,
	}

	return utils.SuccessMessage(response, fiber.StatusOK)
}
