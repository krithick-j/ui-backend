package service

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RecursiveUser struct {
	Name           string         `json:"name"`
	TrackingCenter string         `json:"tracking_center"`
	LeftPoint      string         `json:"left_point"`
	RightPoint     string         `json:"right_point"`
	BV             string         `json:"bv"`
	Left           *RecursiveUser `json:"left"`
	Right          *RecursiveUser `json:"right"`
}

func GetUserByDistId(dist_id string) fiber.Map {

	var user models.User
	var result *gorm.DB

	user, result = repositories.GetUserByID(dist_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": user}
}

func GetUsers() fiber.Map {
	var users []models.User
	var result *gorm.DB
	users, result = repositories.GetAllUsers(users)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": users}
}

func FindRecursiveTC(ruser *RecursiveUser, dist_id string, center_code string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, center_code)
	var nuser *RecursiveUser = new(RecursiveUser)
	nuser.Name = tc.Name
	nuser.TrackingCenter = tc.DistribID + " " + tc.CenterCode
	nuser.LeftPoint = "2500"
	nuser.RightPoint = "3500"
	nuser.BV = "5"
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

func GetTreeUserByDistId(dist_id string) fiber.Map {

	ruser := new(RecursiveUser)
	tc := repositories.GetTrackingCenter(dist_id, "001")
	ruser.Name = tc.Name
	ruser.TrackingCenter = tc.DistribID + " " + tc.CenterCode
	ruser.LeftPoint = "25500"
	ruser.RightPoint = "34500"
	ruser.BV = "10"

	if tc.LeftDistribID != "" {
		FindRecursiveTC(ruser, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID != "" {
		FindRecursiveTC(ruser, tc.RightDistribID, tc.RightPlace, "right")
	}
	return fiber.Map{"data": ruser}
}

func FindNextAvailUserSeq() string {
	last_no := repositories.GetLastId()
	distrib_no, _ := strconv.Atoi(strings.TrimPrefix(last_no, "IN-"))
	return fmt.Sprintf("IN-%05d", distrib_no+1)
}

func FindNextAvailSlot(distrib_id string, side string) string {
	/**
		On a Pyramid network, a reference can only be added either on left or right
		if a person adds thrid person an so on, the actual referree becomes the person below
		the person, if he has an empty slot on the same side
		if not the tree traverse till the botton where it finds an empty slot
		Here we find an empty slot recursively on the same side
		Caution a circular refernce by external db edit may cause an infinite loop
	**/
	var old_distrib_id string
	for {
		old_distrib_id = distrib_id
		distrib_id = repositories.GetSide(distrib_id, side)
		if distrib_id == "" {
			return old_distrib_id
		}
	}
}

func RegisterUser(user_in dto.UserIn) (fiber.Map, error) {
	//TODO: User Table Update
	//Generate Automatically
	distrib_id := FindNextAvailUserSeq()
	user := models.User{
		DistribID: distrib_id,
		Name:      user_in.Name,
		Pass:      fmt.Sprintf("%x", sha256.Sum256([]byte(user_in.Pass))),
	}
	repositories.CreateUser(user)
	tc1 := models.TrackingCenter{
		Name:           user_in.Name,
		DistribID:      distrib_id,
		CenterCode:     "001",
		RefDistribID:   user_in.RefDistID,
		LeftDistribID:  distrib_id,
		LeftPlace:      "002",
		RightDistribID: distrib_id,
		RightPlace:     "003",
	}
	tc2 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		CenterCode: "002",
	}
	tc3 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		CenterCode: "003",
	}
	repositories.CreateTCs([]models.TrackingCenter{tc1, tc2, tc3})
	repositories.UpdateTC(distrib_id, user_in.RefDistID, user_in.RefCenterCode, user_in.Place)
	return fiber.Map{"data": user_in}, nil
}
