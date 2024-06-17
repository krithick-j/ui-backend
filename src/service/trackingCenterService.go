package service

import (
	"fmt"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FindRecursiveTC(ruser *dto.RecursiveUser, dist_id string, place string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, place)
	tcbv, _ := repositories.GetBVforTC(dist_id, place)
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
		FindRecursiveTC(nuser, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID != "" {
		FindRecursiveTC(nuser, tc.RightDistribID, tc.RightPlace, "right")
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}
}

func FindRecursiveTCForRankBv(leftRankBv *float64, rightRankBv *float64, dist_id string, place string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, place)
	tcbv, _ := repositories.GetBVforTC(dist_id, place)
	for _, val := range tcbv {
		if val.Side == "left" {
			*leftRankBv += float64(val.BValue)
		}
		if val.Side == "right" {
			*rightRankBv += float64(val.BValue)
		}
	}

	if tc.LeftDistribID != "" {
		FindRecursiveTCForRankBv(leftRankBv, rightRankBv, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID != "" {
		FindRecursiveTCForRankBv(leftRankBv, rightRankBv, tc.RightDistribID, tc.RightPlace, "right")
	}
}

// only return tracking centers with respect to distrib id
func FindRecursiveTCOnlyDistribId(ruser *dto.RecursiveUser, dist_id string, place string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, place)
	tcbv, _ := repositories.GetBVforTC(dist_id, place)
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
		FindRecursiveTCOnlyDistribId(nuser, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID == tc.DistribID {
		FindRecursiveTCOnlyDistribId(nuser, tc.RightDistribID, tc.RightPlace, "right")
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}
}

func GetTreeUserByDistId(distrib_id string) fiber.Map {

	ruser := new(dto.RecursiveUser)
	tc := repositories.GetTrackingCenter(distrib_id, "001")
	tcbv, _ := repositories.GetBVforTC(distrib_id, "001")
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
		FindRecursiveTC(ruser, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID != "" {
		FindRecursiveTC(ruser, tc.RightDistribID, tc.RightPlace, "right")
	}
	return fiber.Map{"data": ruser}
}

func GetTrackingCentersByDistribId(distrib_id string, isActive bool, fromDate string, toDate string) (fiber.Map, int) {

	tcArr := []dto.TCBv{}
	res, err := repositories.GetAllTrackingCenters(distrib_id)
	if err != nil {
		configs.Log.Errorln("Error on calling GetAllTrackingCenters repositories fn from GetTrackingCentersByDistribId", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	tcOut := dto.TrackingCenterOut{}
	tcOut.DistribId = distrib_id
	for _, place := range res {
		bvres, res := repositories.GetBVforTCOneRowByDateGeneric(distrib_id, place.Place, isActive, fromDate, toDate)
		if res.Error != nil {
			configs.Log.Errorln("Error on calling GetBVforTCOneRow repositories fn from GetTrackingCentersByDistribId", res.Error.Error())
			return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
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
			err := repositories.ActivateTC(distrib_id, placeBv.Place)
			if err != nil {
				fmt.Println("Error on Activating TC")
				configs.Log.Errorln("Error on Activating TC", err.Error())
				tx.Rollback()
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}
			var side string
			if placeBv.Place == "001" {
				side = "bv"
			} else if placeBv.Place == "002" {
				side = "left"
			} else if placeBv.Place == "003" {
				side = "right"
			} else {
				configs.Log.Errorln("Error in updating values of Place")
				return fiber.Map{"error": "Error in updating values of Place"}, fiber.StatusInternalServerError
			}
			BvObj := models.BvTransaction{
				DistribId:    distrib_id,
				Place:        placeBv.Place,
				OrderId:      orderId,
				Date:         time.Now(),
				BvValue:      placeBv.AddBv,
				ActivateDate: time.Now().AddDate(0, 0, 15),
				Side:         side,
				TransType:    "product",
			}
			res := repositories.SaveBvTransaction(BvObj)
			fmt.Printf("res: %v\n", BvObj)
			if res.Error != nil {
				configs.Log.Errorln("Error on saving BvTransaction", res.Error.Error())
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
			place := placeBv.Place
			if side == "bv" {
				side = "left"
			}
			currentTc := repositories.GetTrackingCenter(distrib_id, place)
			for {
				if currentTc.PDistribId == "" {
					configs.Log.Infoln("Breaking from Infinite loop")
					break
				}
				parentTc := repositories.GetTrackingCenter(currentTc.PDistribId, currentTc.PPlace)
				if parentTc.RightDistribID == currentTc.DistribID && parentTc.RightPlace == currentTc.Place {
					side = "right"
				} else {
					//if parentTc.LeftDistribID == currentTc.DistribID && parentTc.LeftPlace == currentTc.Place {
					side = "left"
				}
				BvObj := models.BvTransaction{
					DistribId:    parentTc.DistribID,
					Place:        parentTc.Place,
					OrderId:      orderId,
					Date:         time.Now(),
					BvValue:      placeBv.AddBv,
					ActivateDate: time.Now().AddDate(0, 0, 15),
					Side:         side,
					TransType:    "product",
				}
				if parentTc.IsActive {
					res := repositories.SaveBvTransaction(BvObj)
					if res.Error != nil {
						configs.Log.Errorln("Error on Saving Tc", res.Error.Error())
						tx.Rollback()
						return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
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
