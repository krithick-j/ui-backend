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

func FindRecursiveTC(ruser *RecursiveUser, dist_id string, place string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, place)
	tcbv, _ := repositories.GetBVforTC(dist_id, place)
	var nuser *RecursiveUser = new(RecursiveUser)
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

func GetTreeUserByDistId(distrib_id string) fiber.Map {

	ruser := new(RecursiveUser)
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

	var tracking_centers []models.TrackingCenter
	var result *gorm.DB

	tracking_centers, result = repositories.GetTrackingCenterByDistribId(distrib_id, tracking_centers)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": tracking_centers}, http.StatusOK
}

func UpdateCurrentPlaceValues(distrib_id string, placeBvs []dto.PlaceBv, orderId string) (fiber.Map, int) {
	//Adding Bv Points from the product to the tree
	for _, placeBv := range placeBvs {
		if placeBv.AddBv == 0 {
			//Skip updating for empty values
			continue
		}
		repositories.ActivateTC(distrib_id, placeBv.Place)
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
			return fiber.Map{"error": "Error in updating values of Place"}, http.StatusInternalServerError
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
			tx := models.BvTransaction{
				DisribId:     rdistrib_id,
				Place:        place,
				OrderId:      orderId,
				Date:         time.Now(),
				BvValue:      placeBv.AddBv,
				ActivateDate: time.Now().AddDate(0, 0, 7),
				Side:         nside,
			}
			if currentTc.IsActive {
				repositories.SaveBvTransaction(tx)
			}
			if currentTc.PDistribId == "" {
				break
			}
			rdistrib_id = currentTc.PDistribId
			place = currentTc.PPlace
		}
	}

	return fiber.Map{"data": "Data Successfully Updated"}, http.StatusOK
}
