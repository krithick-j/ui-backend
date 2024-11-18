package service

import (
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

func FindRecursiveTC(ruser *dto.RecursiveUser, dist_id string, place string, side string, tx *gorm.DB) (fiber.Map, int) {
	fmt.Println("----------||||||--------------------")
	tc, err := repositories.GetTrackingCenter(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "FindRecursiveTC", fiber.StatusInternalServerError, tx)
	}

	tcbv, err := repositories.GetBVforTC(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetBVforTC", "FindRecursiveTC", fiber.StatusInternalServerError, tx)
	}

	var nuser *dto.RecursiveUser = new(dto.RecursiveUser)
	nuser.Name = tc.Name
	nuser.TrackingCenter = tc.DistribID + " " + tc.Place
	nuser.IsActive = tc.IsActive
	for _, val := range tcbv {
		if val.Side == "left" {
			nuser.LeftPoint = val.BValue
		}
		if val.Side == "right" {
			nuser.RightPoint = val.BValue
		}
		if val.Side == "bv" {
			nuser.BV = val.BValue
		}
	}
	if tc.LeftDistribID != "" {
		res, status := FindRecursiveTC(nuser, tc.LeftDistribID, tc.LeftPlace, "left", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	if tc.RightDistribID != "" {
		res, status := FindRecursiveTC(nuser, tc.RightDistribID, tc.RightPlace, "right", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}
	return utils.SuccessMessage("sucess", fiber.StatusOK)
}

func FindRecursiveTCForRankBv(leftRankBv *float64, rightRankBv *float64, dist_id string, place string, side string, tx *gorm.DB) (fiber.Map, int) {

	tc, err := repositories.GetTrackingCenter(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "FindRecursiveTCForRankBv", fiber.StatusInternalServerError, tx)
	}
	tcbv, err := repositories.GetBVforTC(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetBVforTC", "FindRecursiveTCForRankBv", fiber.StatusInternalServerError, tx)
	}
	for _, val := range tcbv {
		if val.Side == "left" {
			*leftRankBv += float64(val.BValue)
		}
		if val.Side == "right" {
			*rightRankBv += float64(val.BValue)
		}
	}

	if tc.LeftDistribID != "" {
		res, status := FindRecursiveTCForRankBv(leftRankBv, rightRankBv, tc.LeftDistribID, tc.LeftPlace, "left", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	if tc.RightDistribID != "" {
		res, status := FindRecursiveTCForRankBv(leftRankBv, rightRankBv, tc.RightDistribID, tc.RightPlace, "right", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	return utils.SuccessMessage("Success", fiber.StatusOK)
}

// only return tracking centers with respect to distrib id
func FindRecursiveTCOnlyDistribId(ruser *dto.RecursiveUser, dist_id string, place string, side string, tx *gorm.DB) (fiber.Map, int) {
	fmt.Println(")))))))))))))_-------------------")
	tc, err := repositories.GetTrackingCenter(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "FindRecursiveTCOnlyDistribId", fiber.StatusInternalServerError, tx)
	}
	tcbv, err := repositories.GetBVforTC(dist_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetBVforTC", "FindRecursiveTCOnlyDistribId", fiber.StatusInternalServerError, tx)
	}
	var nuser *dto.RecursiveUser = new(dto.RecursiveUser)
	nuser.Name = tc.Name
	nuser.TrackingCenter = tc.DistribID + " " + tc.Place
	nuser.IsActive = tc.IsActive
	for _, val := range tcbv {
		if val.Side == "left" {
			nuser.LeftPoint = val.BValue
		}
		if val.Side == "right" {
			nuser.RightPoint = val.BValue
		}
		if val.Side == "bv" {
			nuser.BV = val.BValue
		}
	}

	if tc.LeftDistribID == tc.DistribID {
		res, status := FindRecursiveTCOnlyDistribId(nuser, tc.LeftDistribID, tc.LeftPlace, "left", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	if tc.RightDistribID == tc.DistribID {
		res, status := FindRecursiveTCOnlyDistribId(nuser, tc.RightDistribID, tc.RightPlace, "right", tx)
		if status != fiber.StatusOK {
			return res, status
		}
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}
	return utils.SuccessMessage("Sucess", fiber.StatusOK)
}

func GetTreeUserByDistId(distrib_id string, tx *gorm.DB) (fiber.Map, int) {

	ruser := new(dto.RecursiveUser)
	fmt.Println(":::::::::::::::::::::::::::::;")
	tc, err := repositories.GetTrackingCenter(distrib_id, "001", tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "GetTreeUserByDistId", fiber.StatusInternalServerError, tx)
	}
	tcbv, err := repositories.GetBVforTC(distrib_id, "001", tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "GetTreeUserByDistId", fiber.StatusInternalServerError, tx)
	}

	ruser.Name = tc.Name
	ruser.TrackingCenter = tc.DistribID + " " + tc.Place
	ruser.IsActive = tc.IsActive
	for _, val := range tcbv {
		if val.Side == "left" {
			ruser.LeftPoint = val.BValue
		}
		if val.Side == "right" {
			ruser.RightPoint = val.BValue
		}
		if val.Side == "bv" {
			ruser.BV = val.BValue
		}
	}

	if tc.LeftDistribID != "" {
		FindRecursiveTC(ruser, tc.LeftDistribID, tc.LeftPlace, "left", tx)
	}
	if tc.RightDistribID != "" {
		FindRecursiveTC(ruser, tc.RightDistribID, tc.RightPlace, "right", tx)
	}
	return utils.SuccessMessage(ruser, fiber.StatusOK)
}

func GetTrackingCentersByDistribId(distrib_id string, isActive bool, fromDate string, toDate string, tx *gorm.DB) (fiber.Map, int) {

	tcArr := []dto.TCBv{}
	res, err := repositories.GetAllTrackingCenters(distrib_id, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllTrackingCenters", "GetTrackingCentersByDistribId", fiber.StatusInternalServerError, tx)
	}

	tcOut := dto.TrackingCenterOut{}
	tcOut.DistribId = distrib_id
	for _, place := range res {
		bvres, err := repositories.GetBVforTCOneRowByDateGeneric(distrib_id, place.Place, isActive, fromDate, toDate, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetAllTrackingCenters", "GetTrackingCentersByDistribId", fiber.StatusInternalServerError, tx)
		}
		obj := dto.TCBv{
			Place:   place.Place,
			LPoint:  bvres.LValue,
			RPoint:  bvres.RValue,
			BvPoint: bvres.BValue,
		}
		tcArr = append(tcArr, obj)
	}

	tcOut.Tc = tcArr

	return fiber.Map{"data": tcOut}, fiber.StatusOK
}

func GetTrackingCentersByDistribIdForCheque(distrib_id string, fromDate string, toDate string, tx *gorm.DB) (fiber.Map, int) {

	tcArr := []dto.TCBv{}
	res, err := repositories.GetAllTrackingCenters(distrib_id, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllTrackingCenters", "GetTrackingCentersByDistribId", fiber.StatusInternalServerError, tx)
	}

	tcOut := dto.TrackingCenterOut{}
	tcOut.DistribId = distrib_id
	for _, place := range res {
		bvres, err := repositories.GetBVforTCOneRowByDateForCheque(distrib_id, place.Place, fromDate, toDate, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetAllTrackingCenters", "GetTrackingCentersByDistribId", fiber.StatusInternalServerError, tx)
		}
		obj := dto.TCBv{
			Place:   place.Place,
			LPoint:  bvres.LValue,
			RPoint:  bvres.RValue,
			BvPoint: bvres.BValue,
		}
		tcArr = append(tcArr, obj)
	}

	tcOut.Tc = tcArr

	return fiber.Map{"data": tcOut}, fiber.StatusOK
}

func UpdateCurrentPlaceValues(distrib_id string, placeBvs []dto.PlaceBv, orderId string, tx *gorm.DB, totalBv float64) (fiber.Map, int) {

	sum := 0.0
	configs.Log.Infoln("place bv array===>", placeBvs)
	for _, placeBv := range placeBvs {
		sum += placeBv.AddBv
	}

	configs.Log.Infoln("sum is ", sum, "total bv is", totalBv)
	//Validation--> Sum of bv should match the totalValueType
	if sum == totalBv {
		//Adding Bv Points from the product to the tree
		for _, placeBv := range placeBvs {
			if placeBv.AddBv == 0 {
				//Skip updating for empty values
				continue
			}
			err := repositories.ActivateTC(distrib_id, placeBv.Place, tx)
			if err != nil {
				return utils.NotNilErrorMessage(err, "ActivateTC", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
			}
			trackingCenter, err := repositories.GetTrackingCenter(distrib_id, placeBv.Place, tx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return utils.NotNilErrorMessage(err, "GetTrackingCenter", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
			}
			parentTrackingCenter, err := repositories.GetTrackingCenter(trackingCenter.PDistribId, trackingCenter.PPlace, tx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return utils.NotNilErrorMessage(err, "GetTrackingCenter", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
			}
			var side string
			if parentTrackingCenter.LeftDistribID == trackingCenter.DistribID && parentTrackingCenter.LeftPlace == trackingCenter.Place {
				side = "left"
			} else {
				side = "right"
			}

			//hardcoded for 001
			if placeBv.Place == "001" {
				side = "bv"
			}

			// Example current date
			currentDate := time.Now()

			// Calculate the next Friday
			daysUntilFriday := (5 - int(currentDate.Weekday()) + 7) % 7
			nextFriday := currentDate.AddDate(0, 0, daysUntilFriday)

			// Add 14 days to get the desired ActivateDate
			activateDate := nextFriday.AddDate(0, 0, 14)
			BvObj := models.BvTransaction{
				DistribId:    distrib_id,
				Place:        placeBv.Place,
				OrderId:      orderId,
				Date:         time.Now(),
				BvValue:      placeBv.AddBv,
				ActivateDate: activateDate,
				Side:         side,
				TransType:    "product",
			}
			err = repositories.SaveBvTransaction(BvObj, tx)
			if err != nil {
				return utils.NotNilErrorMessage(err, "SaveBvTransaction", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
			}
			place := placeBv.Place
			if side == "bv" {
				side = "left"
			}
			currentTc, err := repositories.GetTrackingCenter(distrib_id, place, tx)
			if err != nil {
				return utils.NotNilErrorMessage(err, "SaveBvTransaction", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
			}
			for {
				if currentTc.PDistribId == "" {
					configs.Log.Infoln("Breaking from Infinite loop")
					break
				}
				parentTc, err := repositories.GetTrackingCenter(currentTc.PDistribId, currentTc.PPlace, tx)
				if err != nil {
					return utils.NotNilErrorMessage(err, "GetTrackingCenter", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
				}
				if parentTc.RightDistribID == currentTc.DistribID && parentTc.RightPlace == currentTc.Place {
					side = "right"
				} else {
					//if parentTc.LeftDistribID == currentTc.DistribID && parentTc.LeftPlace == currentTc.Place {
					side = "left"
				}
				// Example current date
				currentDate := time.Now()

				// Calculate the next Friday
				daysUntilFriday := (5 - int(currentDate.Weekday()) + 7) % 7
				nextFriday := currentDate.AddDate(0, 0, daysUntilFriday)

				BvObj := models.BvTransaction{
					DistribId:    parentTc.DistribID,
					Place:        parentTc.Place,
					OrderId:      orderId,
					Date:         time.Now(),
					BvValue:      placeBv.AddBv,
					ActivateDate: nextFriday,
					Side:         side,
					TransType:    "product",
				}
				if parentTc.IsActive {
					err := repositories.SaveBvTransaction(BvObj, tx)
					if err != nil {
						return utils.NotNilErrorMessage(err, "SaveBvTransaction", "UpdateCurrentPlaceValues", fiber.StatusInternalServerError, tx)
					}

				}
				currentTc = parentTc
			}
		}
		return fiber.Map{"data": "Data Successfully Updated"}, fiber.StatusOK
	}
	configs.Log.Errorln("total value and Total Bv in distribution table does not match!!Check the input values")
	return fiber.Map{"error": "total value and Total Bv in distribution table does not match!!Check the input values"}, fiber.StatusInternalServerError
}

// available tc to buy is found using this function
func GetAvailableAddTc(distribId string, tx *gorm.DB) (fiber.Map, int) {
	// get bv sum both active and not active
	bvSum, err := repositories.GetBvSumByDistribId(distribId, "product", tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetBvSumByDistribId", "AddTc", fiber.StatusInternalServerError, tx)
	}
	str := fmt.Sprintf("%v", bvSum)
	if bvSum < 1000 {
		return utils.SuccessMessage(0, fiber.StatusForbidden)
	} else if bvSum > 1000 && bvSum < 10000 {
		return utils.SuccessMessage(int(str[0]-'0'), fiber.StatusOK) // int('7' - '0')  '7' is 55 in ASCII, '0' is 48, so 55 - 48 = 7
	} else if bvSum > 10000 && bvSum < 100000 {
		return utils.SuccessMessage(int(str[0]-'0')*10+int(str[1]-'0'), fiber.StatusOK) // 38500 = 38
	} else if bvSum > 100000 && bvSum < 1000000 {
		return utils.SuccessMessage(int(str[0]-'0')*100+int(str[1]-'0')*10+int(str[2]-'0'), fiber.StatusOK) // 388500 = 388
	} else {
		return utils.CommonMessage("Something went wrong!", fiber.StatusInternalServerError, tx)
	}
}

func AddTc(payload dto.AddTc, tx *gorm.DB) (fiber.Map, int) {

	//get tc active and not active length
	tcArr, err := repositories.GetAllTrackingCenters(payload.DistribID, tx)
	if err != nil {
		configs.Log.Errorln("Error on calling GetAllTrackingCenters repositories fn from AddTc fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	res, status := GetAvailableAddTc(payload.DistribID, tx)
	if status != fiber.StatusOK {
		return res, status
	}
	bvSumFinal := res["data"].(int)
	//get current tc maximum number
	res, status = FindNextTcNumber(payload.DistribID, tx)
	if status != fiber.StatusOK {
		return res, status
	}
	place := res["data"].(string)
	// Get the next tracking center number
	tcLen := len(tcArr)

	if tcLen >= bvSumFinal {

		configs.Log.Errorln("Error on checking tcLend and bvSumFinal from AddTc service")
		return fiber.Map{"data": "Tracking Center cannot be created"}, fiber.StatusForbidden
	} else if tcLen < bvSumFinal {

		userData, err := repositories.GetUserByID(payload.DistribID, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetUserByID", "AddTc", fiber.StatusInternalServerError, tx)
		}

		//if not empty don't overwrite but find next available free slot
		parent_distrib_id, parent_ref_place, err := FindNextAvailSlot(payload.PlacementDistribId, payload.PlacementPlace, payload.PlacementSide, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "FindNextAvailSlot", "AddTc", fiber.StatusInternalServerError, tx)
		}
		newTc := models.TrackingCenter{
			Name:       userData.Name,
			DistribID:  payload.DistribID,
			Place:      place,
			PDistribId: parent_distrib_id,
			PPlace:     parent_ref_place,
			IsActive:   false,
		}

		err = repositories.CreateTCs(newTc, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "CreateTCs", "AddTc", fiber.StatusInternalServerError, tx)
		}

		err = repositories.UpdateTC(tx, payload.DistribID, place, parent_distrib_id, parent_ref_place, payload.PlacementSide)
		if err != nil {
			return utils.NotNilErrorMessage(err, "UpdateTC", "AddTc", fiber.StatusInternalServerError, tx)
		}
	}
	return utils.SuccessMessage(place, fiber.StatusCreated)
}

func FindNextTcNumber(distribId string, tx *gorm.DB) (fiber.Map, int) {
	tcArr, err := repositories.GetAllTrackingCenters(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllTrackingCenters", "FindNextTcNumber", fiber.StatusInternalServerError, tx)
	}

	// Get the next tracking center number
	nextTcNumber := len(tcArr) + 1 // 1 for next tc number

	result := fmt.Sprintf("%03d", nextTcNumber)
	return utils.SuccessMessage(result, fiber.StatusOK)
	// Format the number to a 3-digit string with leading zeros
}
