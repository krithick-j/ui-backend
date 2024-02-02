package service

import (
	"fmt"
	"strconv"
	"strings"
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

func FindRecursiveFind(ruser *RecursiveUser, dist_id string, side string) {
	var user models.User
	var result *gorm.DB
	user, result = repositories.GetUserByID(dist_id, user)
	var nuser *RecursiveUser = new(RecursiveUser)
	nuser.Name = user.Name
	nuser.TrackingCenter = user.DistID + "001"
	nuser.LeftPoint = "24500"
	nuser.RightPoint = "23500"
	nuser.BV = "50"

	if result.Error == gorm.ErrRecordNotFound {
		panic("Not Found")
	}
	if result.Error != nil {
		panic("Some Error")
	}
	if user.Lside != "" {
		FindRecursiveFind(nuser, user.Lside, "left")
	}
	if user.Rside != "" {
		FindRecursiveFind(nuser, user.Rside, "right")
	}
	if side == "left" {
		ruser.Left = nuser
	} else {
		ruser.Right = nuser
	}

}

func GetTreeUserByDistId(dist_id string) fiber.Map {

	var user models.User
	var result *gorm.DB
	data := new(RecursiveUser)
	user, result = repositories.GetUserByID(dist_id, user)
	data.Name = user.Name
	if user.Lside != "" {
		FindRecursiveFind(data, user.Lside, "left")
	}
	if user.Rside != "" {
		FindRecursiveFind(data, user.Rside, "right")
	}
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": data}
}

func FindNextAvailUserSeq() string {
	last_no := repositories.GetLastId()
	distrib_no, _ := strconv.Atoi(strings.TrimPrefix(last_no, "IN-"))
	return fmt.Sprintf("IN-%03d", distrib_no+1)
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

func RegisterUser(regDetails models.User) (fiber.Map, error) {
	var side string
	if regDetails.Place == "L" {
		side = "lside"
	} else {
		side = "rside"
	}
	regDetails.DistID = FindNextAvailUserSeq()
	regDetails.RefDistID = FindNextAvailSlot(regDetails.RefDistID, side)
	err := repositories.CreateUser(regDetails)
	if err != nil {
		return nil, err
	}

	return fiber.Map{"data": regDetails}, nil
}
