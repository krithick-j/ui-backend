package service

import (
	"crypto/sha256"
	"fmt"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get Total Cheque Value By Distrib ID
func TotalChequeValueByDistribId(TakeChequeIn dto.CheckoutIn, tx *gorm.DB) (fiber.Map, int) {

	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000

	//Get Tracking Center with is_active = 1
	trackingCentersActiveValue, status := GetTrackingCentersByDistribIdForCheque(TakeChequeIn.DistribId, "", "", tx)
	if status != fiber.StatusOK {
		return trackingCentersActiveValue, status
	}
	configs.Log.Infoln("Get tracking centers completed")

	//use type assertion to map the variables
	trackingCentersAllArray, ok := trackingCentersActiveValue["data"].(dto.TrackingCenterOut)
	if !ok {
		return fiber.Map{"data": "something went wrong in type assertion tracking centers", "status": ok}, fiber.StatusInternalServerError
	}

	var chequeFrequencyArr []dto.ChequeFrequency
	//Formula and addition
	for _, tc := range trackingCentersAllArray.Tc {

		ChequeFrequency := utils.NCheckoutPossible(tc.LPoint, tc.RPoint, CHEQUE_DRAW_VALUE)
		//Adding is_active = 0 values
		cfObj := dto.ChequeFrequency{
			Tc:        tc.Place,
			Frequency: int(ChequeFrequency),
		}
		chequeFrequencyArr = append(chequeFrequencyArr, cfObj)
	}

	//Dto Out
	response := dto.TakeChequeOut{
		AvailableChequeCountArr: chequeFrequencyArr,
	}

	return utils.SuccessMessage(response, fiber.StatusOK)
}

func GetFrequencyAmount(payload dto.FrequencyForTc, tx *gorm.DB) (fiber.Map, int) {

	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000
	COUNT := 2                                                //Left and Right inside the tracking center
	var cpaAmount float64
	var epAmount float64
	//Get Rank value
	rank, err := repositories.GetCurrentRankValueByDistribId(payload.DistribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetRankValueByDistribId", "TotalChequeValueByDistribId", fiber.StatusInternalServerError, tx)
	}
	count, err := repositories.GetCheckoutFrequency(payload.DistribId, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.NotNilErrorMessage(err, "GetCheckoutFrequency", "GetFrequencyAmount", fiber.StatusInternalServerError, tx)
	}
	freqAmount := float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank

	if payload.Frequency <= 0 {
		epAmount = 0
		cpaAmount = 0
		freqAmount = 0
	} else {
		newChequeFrequency := payload.Frequency + count
		//Formula for ep multiplier
		var epMultiplier int
		for i := count; i <= newChequeFrequency; i++ {
			if (i+1)%5 == 0 {
				fmt.Println("i--->", i)
				epMultiplier += 1
			}
		}
		cpaMultiplier := payload.Frequency - epMultiplier
		fmt.Println("cpa multiplier--->", cpaMultiplier)
		fmt.Println("ep multiplier--->", epMultiplier)
		fmt.Println("freqamount--->", freqAmount)
		cpaAmount = float64(cpaMultiplier) * freqAmount
		epAmount = float64(epMultiplier) * freqAmount
	}

	response := map[string]float64{
		"ep_amount":  epAmount,
		"cpa_amount": cpaAmount,
	}
	return utils.SuccessMessage(response, fiber.StatusOK)
}

func GetValuesForCpa(distribId string, tx *gorm.DB) (fiber.Map, int) {

	totalCpaBalance, err := repositories.GetCpaBalance(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCpaBalance", "GetValuesForCpa", fiber.StatusInternalServerError, tx)
	}

	totalAvailableCpaBalance, err := repositories.GetAvailableCpaBalance(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCpaBalance", "GetValuesForCpa", fiber.StatusInternalServerError, tx)
	}

	totalDcBalance, err := repositories.GetDirectCommissionAllValueByDistribId(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetDirectCommissionValueByDistribId", "GetValuesForCpa", fiber.StatusInternalServerError, tx)
	}

	//Active Points start.....
	totalAvailalbeDcBalance, err := repositories.GetDirectCommissionActiveValueByDistribId(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetDirectCommissionActiveValueByDistribId", "GetValuesForCpa", fiber.StatusInternalServerError, tx)
	}

	//totalAvailablePoints is only is_active
	totalCpaDcBalance := totalCpaBalance + totalCpaBalance
	totalAvailalbeCpaDcBalance := totalAvailableCpaBalance + totalAvailableCpaBalance

	type Response struct {
		DirectComissionTotalBalance          float64 `json:"dc_total_balance"`
		DirectComissionTotalAvailableBalance float64 `json:"dc_total_avail_balance"`
		CpaTotalBalance                      float64 `json:"cpa_total_balance"`
		CpaTotalAvailableBalance             float64 `json:"bv_total_avail_balance"`
		TotalCpaDcBalance                    float64 `json:"total_cpa_dc_balance"`
		TotalCpaDcAvailableBalance           float64 `json:"total_cpa_dc_avail_balance"`
	}

	response := Response{
		DirectComissionTotalBalance:          totalDcBalance,
		DirectComissionTotalAvailableBalance: totalAvailalbeDcBalance,
		CpaTotalBalance:                      totalCpaBalance,
		CpaTotalAvailableBalance:             totalAvailableCpaBalance,
		TotalCpaDcBalance:                    totalCpaDcBalance,
		TotalCpaDcAvailableBalance:           totalAvailalbeCpaDcBalance,
	}
	return utils.SuccessMessage(response, fiber.StatusOK)
}

func TakeChequeByDistribIdAndPlace(TakeChequeIn dto.TakeChequeIn, tx *gorm.DB) (fiber.Map, int) {

	checkoutId := GenerateUniqueHexCode(10)
	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000
	COUNT := 2                                                //Left and Right inside the tracking center

	for i := 1; i <= TakeChequeIn.ChequeCount; i++ {

		tcBv, err := repositories.GetBVforTCOneRow(TakeChequeIn.DistribId, TakeChequeIn.Place, tx)
		if err != nil {
			configs.Log.Errorln("Error on calling GetBVforTCOneRow from TakeChequeByDistribIdAndPlace", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		totalChequeCount := utils.NCheckoutPossible(tcBv.LValue, tcBv.RValue, CHEQUE_DRAW_VALUE)
		if totalChequeCount < float64(TakeChequeIn.ChequeCount) {
			tx.Rollback()
			configs.Log.Errorln("Error from totalChequeCount validation. Cheque cannot be taken because of insufficient BV values")
			return fiber.Map{"data": "cheque unsuccessfull!"}, fiber.StatusForbidden
		}
		if tcBv.LValue < CHEQUE_DRAW_VALUE || tcBv.RValue < CHEQUE_DRAW_VALUE {
			tx.Rollback()
			configs.Log.Errorln("Error from chequeDrawValue validation. Cheque cannot be taken because of insufficient BV values: ")
			return fiber.Map{"data": "cheque unsuccessfull!"}, fiber.StatusForbidden
		}

		count, err := repositories.GetCheckoutFrequency(TakeChequeIn.DistribId, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.NotNilErrorMessage(err, "CreateCheckoutFrequency", "TakeChequeByDistribIdAndPlace", fiber.StatusInternalServerError, tx)
		}
		if err == gorm.ErrRecordNotFound {
			if err = repositories.CreateCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, tx); err != nil {
				return utils.NotNilErrorMessage(err, "CreateCheckoutFrequency", "TakeChequeByDistribIdAndPlace", fiber.StatusInternalServerError, tx)
			}
		}

		rank, err := repositories.GetCurrentRankValueByDistribId(TakeChequeIn.DistribId, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetRankValueByDistribId", "handleChequeBvTransaction", fiber.StatusInternalServerError, tx)
		}

		oneFrequencyAmount := float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank

		//This validation mostly won't be needed here because this is already done in get total cheque count frequency. Anyway it is used for extra validation
		//Suppose 4 cheques are taken, 3 will be added to cpa, fourth cheque will be added to ep.
		//Here count is supposed to check with 3, because the loop is handling the fourth cheque.
		if count%3 == 0 && count > 0 {
			//take this cheque, this will not add in the cpa balance
			res, status := handleEpAmount(TakeChequeIn.DistribId, checkoutId, oneFrequencyAmount, tx)
			if status != fiber.StatusOK {
				return res, status
			}
		} else {
			res, status := handleCpaTransaction(TakeChequeIn.DistribId, oneFrequencyAmount, checkoutId, tx)
			if status != fiber.StatusOK {
				return res, status
			}
		}
		res, status := handleChequeBvTransaction(TakeChequeIn, checkoutId, CHEQUE_DRAW_VALUE, tx)
		if status != fiber.StatusOK {
			return res, status
		}
		err = repositories.IncrementCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, count, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "IncrementCheckoutFrequency", "TakeChequeByDistribIdAndPlace", fiber.StatusInternalServerError, tx)
		}
	}
	return utils.SuccessMessage("Cheque taken successful", fiber.StatusOK)
}

func handleChequeBvTransaction(TakeChequeIn dto.TakeChequeIn, checkoutId string, CHEQUE_DRAW_VALUE int, tx *gorm.DB) (fiber.Map, int) {
	LeftInsideTcObj := models.BvTransaction{
		DistribId:    TakeChequeIn.DistribId,
		Place:        TakeChequeIn.Place,
		OrderId:      checkoutId,
		Date:         time.Now(),
		BvValue:      -float64(CHEQUE_DRAW_VALUE),
		ActivateDate: time.Now(),
		Side:         "left",
		TransType:    "cheque",
	}
	RightInsidetCObj := models.BvTransaction{
		DistribId:    TakeChequeIn.DistribId,
		Place:        TakeChequeIn.Place,
		OrderId:      checkoutId,
		Date:         time.Now(),
		BvValue:      -float64(CHEQUE_DRAW_VALUE),
		ActivateDate: time.Now(),
		Side:         "right",
		TransType:    "cheque",
	}

	err := repositories.SaveBvTransaction(LeftInsideTcObj, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveBvTransaction", "handleChequeBvTransaction", fiber.StatusInternalServerError, tx)
	}
	err = repositories.SaveBvTransaction(RightInsidetCObj, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveBvTransaction", "handleChequeBvTransaction", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage("Success on Cheque Bv Transaction!", fiber.StatusOK)
}

func handleEpAmount(distribId, reference string, value float64, tx *gorm.DB) (fiber.Map, int) {

	err := repositories.SaveEpTx(distribId, reference, value, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveBvTransaction", "handleChequeBvTransaction", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage("Success on handling Ep Amount!", fiber.StatusOK)
}

func handleCpaTransaction(distribID string, cpaAmount float64, checkoutId string, tx *gorm.DB) (fiber.Map, int) {

	cpaTransactionObj := models.CpaTransaction{
		DistribId:    distribID,
		Reference:    checkoutId,
		ActivateDate: time.Now(),
		IsActive:     true,
		Amount:       cpaAmount,
	}
	err := repositories.SaveCpaTransaction(cpaTransactionObj, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveCpaTransaction", "TakeChequeByDistribIdAndPlace", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage("Cpa Transaction success!", fiber.StatusOK)
}

func ChangeChequePin(payload dto.ChequePinIn, tx *gorm.DB) (fiber.Map, int) {
	currentPinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.CurrentPin)))

	currentPinHashFromDB, err := repositories.GetChequePinByDistribID(payload.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChangeChequePin service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	if currentPinHashFromDB != currentPinHashFromPayload {
		tx.Rollback()
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	err = repositories.ChangeCpaPin(payload.DistribId, fmt.Sprintf("%x", sha256.Sum256([]byte(payload.NewPin))), tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling ChangeCpaPin repositories fn from ChangeChequePin service fn")
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Cpa pin changed Successfully"}, fiber.StatusOK
}

func ChequeLogin(payload dto.ChequeLogin, tx *gorm.DB) (fiber.Map, int) {
	PinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.Pin)))

	currentPinHashFromDB, err := repositories.GetChequePinByDistribID(payload.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChequeLogin service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError

	}

	if currentPinHashFromDB != PinHashFromPayload {
		tx.Rollback()
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	return fiber.Map{"data": "Cpa pin login successfull"}, fiber.StatusOK
}
