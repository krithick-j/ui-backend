package service

import (
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
		return utils.NotNilErrorMessage(err, "GetReferredUsersByDistribId", "GetGroupRspByDistribId", fiber.StatusInternalServerError, tx)
	}

	for _, referredDistribId := range referredDistribIds {

		totalRspForOneDistrib, err := repositories.GetPersonalRspSumByDistribId(referredDistribId, tx)
		if err != nil {
			utils.NotNilErrorMessage(err, "GetPersonalRspSumByDistribIdt", "GetGroupRspByDistribId", fiber.StatusInternalServerError, tx)
		}

		rspSum += totalRspForOneDistrib

		recursiveRspSum, status := GetGroupRspByDistribId(referredDistribId, tx)
		if status != fiber.StatusOK {
			return recursiveRspSum, status
		}
		rspSum += recursiveRspSum["data"].(float64)
	}

	return fiber.Map{"data": rspSum}, fiber.StatusOK
}

// This function will give a tree to traverse down like a chain of referrals
// Suppose IN-00001 refers IN-00002, and IN-00002 refers IN-00003, then the array returns [IN-00002,IN-00003] if distribId is IN-00001.
// Returns [IN-00003] if input is IN-00002
func GetGroupPerformanceByDistribId(distribId string, tx *gorm.DB) (fiber.Map, int) {
	// Get the referral chain for the given distribId
	res, status := GetReferralChainByDistribId(distribId, tx)
	if status != fiber.StatusOK {
		return res, status
	}

	referredDistribIds := res["data"].([]string)

	// Get the current rank of the original distributor
	currentRankDistribId, err := repositories.GetCurrentRankValueByDistribId(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCurrentRankValueByDistribId", "GetGroupPerformanceByDistribId", fiber.StatusInternalServerError, tx)
	}

	// Count the number of referred distributors with rank >= the original distributor
	var count int64
	if len(referredDistribIds) > 0 {
		// Use a single query to get the ranks of all referred distributors
		referredRanks, err := repositories.GetCurrentRankArrByDistribId(referredDistribIds, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetCurrentRankArrByDistribId", "GetGroupPerformanceByDistribId", fiber.StatusInternalServerError, tx)
		}

		// Count the number of ranks that are >= currentRankDistribId
		for _, rank := range referredRanks {
			if rank >= currentRankDistribId {
				count++
			}
		}
	}

	// Return the count as the performance metric
	return utils.SuccessMessage(count, fiber.StatusOK)
}

func GetDirectBvByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {

	directBv, err := repositories.GetDirectBvByDistribID(DistribID, tx)
	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": directBv}, fiber.StatusOK
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetReferredUsersByDistribId", "GetGroupRspByDistribId", fiber.StatusInternalServerError, tx)
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
