package service

import (
	"fmt"
	"net/http"
	"time"
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

func GetTrackingCentersByDistribId(distrib_id string) (fiber.Map, int) {

	// places := [3]string{"001", "002", "003"}
	// var trackingCenters [][]models.TCBv
	tcArr := []dto.TCBv{}
	res, err := repositories.GetAllTrackingCenters(distrib_id)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	tcOut := dto.TrackingCenterOut{}
	tcOut.DistribId = distrib_id
	for _, place := range res {
		// fmt.Println(r.Place)
		bvres, _ := repositories.GetBVforTCOneRow(distrib_id, place.Place)
		obj := dto.TCBv{
			Place:   place.Place,
			LPoint:  bvres.LValue,
			RPoint:  bvres.RValue,
			BvPoint: bvres.BValue,
		}
		tcArr = append(tcArr, obj)
	}

	tcOut.Tc = tcArr

	return fiber.Map{"data": tcOut}, http.StatusOK
}

func UpdateCurrentPlaceValues(distrib_id string, placeBvs []dto.PlaceBv, orderId string, tx *gorm.DB, totalBv float64) (fiber.Map, int) {

	sum := 0.0
	for _, placeBv := range placeBvs {
		sum += placeBv.AddBv
	}

	fmt.Println("sum", sum, "placebv", placeBvs, "total bv", totalBv)
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
				tx.Rollback()
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}
			var side string
			bvretain_flag := true
			if placeBv.Place == "001" {
				side = "bv"

			} else if placeBv.Place == "002" {
				side = "left"
			} else if placeBv.Place == "003" {
				side = "right"
			} else {
				fmt.Println("Error in updating values of Place")
				return fiber.Map{"error": "Error in updating values of Place"}, fiber.StatusInternalServerError
			}
			place := placeBv.Place
			rdistrib_id := distrib_id
			nside := side
			for {
				currentTc := repositories.GetTrackingCenter(rdistrib_id, place)
				parentPlace := repositories.GetTrackingCenter(currentTc.DistribID, currentTc.PPlace)

				//Deciding left or right
				//Dont change if it is bv and its first time
				if !(place == "001" && bvretain_flag) {
					if parentPlace.RightDistribID == rdistrib_id && parentPlace.RightPlace == currentTc.Place {
						nside = "right"
					}
					if parentPlace.LeftDistribID == rdistrib_id && parentPlace.LeftPlace == currentTc.Place {
						nside = "left"
					}
				}
				bvretain_flag = false
				BvObj := models.BvTransaction{
					DistribId:    rdistrib_id,
					Place:        place,
					OrderId:      orderId,
					Date:         time.Now(),
					BvValue:      placeBv.AddBv,
					ActivateDate: time.Now().AddDate(0, 0, 7),
					Side:         nside,
					TransType:    "Product",
				}
				if currentTc.IsActive {
					res := repositories.SaveBvTransaction(BvObj)
					if res.Error != nil {
						fmt.Println("Error on Saving Tc", res.Error.Error())
						tx.Rollback()
						return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
					}

				}
				if currentTc.PDistribId == "" {
					break
				}
				rdistrib_id = currentTc.PDistribId
				place = currentTc.PPlace
			}
		}
		return fiber.Map{"data": "Data Successfully Updated"}, fiber.StatusOK
	}
	return fiber.Map{"error": "total value and Total Bv in distribution table does not match!!Check the input values"}, fiber.StatusInternalServerError
}
