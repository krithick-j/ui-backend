package service

import (
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
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
	var personalRspMax, groupRspMax, directBvMax, stepMax, groupPerformanceMax int
	rankValue, err := repositories.GetCurrentRankValueByDistribId(DistribID, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCurrentRankValueByDistribId", "GetRspValuesByDistribID", fiber.StatusInternalServerError, tx)
	}

	rankDetails, err := repositories.GetRankDetailsByDistribId(rankValue, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetRankDetailsByDistribId", "GetRspValuesByDistribID", fiber.StatusInternalServerError, tx)
	}
	for _, rankRow := range rankDetails {
		switch rankRow.Type {
		case "PRSP":
			personalRspMax = rankRow.Target
		case "GRSP":
			groupRspMax = rankRow.Target
		case "DRBV":
			directBvMax = rankRow.Target
		case "STEP":
			stepMax = rankRow.Target
		case "GPRF":
			groupPerformanceMax = rankRow.Target
		}
	}
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

	groupPerformanceResult, status := GetGroupPerformanceByDistribId(DistribID, tx)

	if status != fiber.StatusOK {
		return stepResult, status
	}
	grpPerformance := groupPerformanceResult["data"]

	rspUpdateTime, err := repositories.GetRspTime(tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetRspTime", "GetRspValuesByDistribID", fiber.StatusInternalServerError, tx)
	}
	fmtUpdateTime := utils.FormatTimeByLocation(rspUpdateTime, "Asia/Kolkata", "2 Jan 2006 15:04")

	response := fiber.Map{
		"data": fiber.Map{
			"direct_bv": fiber.Map{
				"value":     directBv,
				"max_value": directBvMax,
			},
			"personal_rsp": fiber.Map{
				"value":     personalRsp,
				"max_value": personalRspMax,
			},
			"group_rsp": fiber.Map{
				"value":     groupRsp,
				"max_value": groupRspMax,
			},
			"step": fiber.Map{
				"value":     step,
				"max_value": stepMax,
			},
			"group_performance": fiber.Map{
				"value":     grpPerformance,
				"max_value": groupPerformanceMax,
			},
		},
		"updated_time": fmtUpdateTime,
	}
	return utils.SuccessMessage(response, fiber.StatusOK)
}

func SaveCpaICoupon(payload dto.TakeCpaAmount, tx *gorm.DB) (fiber.Map, int) {

	reference := GenerateUniqueHexCode(10)
	var icouponBalance float64
	cpaBalance, err := repositories.GetAvailableCpaBalance(payload.DistribID, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCpaBalance", "SaveCpaICoupon", fiber.StatusInternalServerError, tx)
	}

	totalAvailalbeDcBalance, err := repositories.GetDirectCommissionActiveValueByDistribId(payload.DistribID, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetDirectCommissionActiveValueByDistribId", "GetValuesForCpa", fiber.StatusInternalServerError, tx)
	}

	for _, icoupon := range payload.Coupons {
		icouponBalance = icouponBalance + (icoupon.Value * float64(icoupon.Quantity))
	}

	if payload.CpaAmount > cpaBalance {
		return utils.CommonMessage("Insufficient cpa amount", fiber.StatusInternalServerError, tx)
	}

	if payload.DcAmount > totalAvailalbeDcBalance {
		return utils.CommonMessage("Insufficient dc amount", fiber.StatusInternalServerError, tx)
	}

	totalBalance := payload.CpaAmount + payload.DcAmount
	if totalBalance != icouponBalance {
		return utils.CommonMessage("Icoupon Balance and total balance does not match", fiber.StatusInternalServerError, tx)
	}
	activateDate := time.Now()
	cpaObj := models.CpaTransaction{
		DistribId:    payload.DistribID,
		Reference:    reference,
		ActivateDate: activateDate,
		IsActive:     true,
		Amount:       -payload.CpaAmount, //cpa value detected
	}
	err = repositories.SaveCpaTransaction(cpaObj, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveCpaTransaction", "SaveCpaICoupon", fiber.StatusInternalServerError, tx)
	}
	if payload.DcAmount != 0 {
		dcObj := models.DirectCommissionTransaction{
			DistribId:    "",
			Value:        -payload.DcAmount,
			Reference:    reference,
			RefDistribId: payload.DistribID, //This will be used to calculate the total dc amount and used to reduce the amount
			ActivateDate: activateDate,
			// ExpiryDate:   activateDate.AddDate(0, 6, 0),
			IsActive:     true,
		}

		err = repositories.SaveDirectCommissionTransaction(dcObj, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "SaveDirectCommissionTransaction", "SaveCpaICoupon", fiber.StatusInternalServerError, tx)
		}
	}

	res, status := AddICoupon(payload.ICouponIn, payload.DistribID, activateDate.AddDate(0, 6, 0), tx)
	if status != fiber.StatusCreated {
		return res, status
	}
	return utils.SuccessMessage(res["data"], status)
}
