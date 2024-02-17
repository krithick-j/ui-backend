package service

/**

A center code is a tc code which will have three values of 001,002,003 and will be created
by default. The same center code if referred in other places called place

**/

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RecursiveUser struct {
	Name           string         `json:"name"`
	TrackingCenter string         `json:"tracking_center"`
	LeftPoint      string         `json:"left_point"`
	RightPoint     string         `json:"right_point"`
	BV             int            `json:"bv"`
	Left           *RecursiveUser `json:"left"`
	Right          *RecursiveUser `json:"right"`
}

func LoginUser(username string, password string) (fiber.Map, int) {
	pass := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	res, err := repositories.AuthUser(username, pass)
	if err != nil {
		return fiber.Map{"err": err.Error()}, http.StatusUnauthorized
	}
	claims := jwt.MapClaims{
		"name":  res.Name,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	tokenstring, err := token.SignedString([]byte("secret"))
	if err != nil {
		return fiber.Map{"err": err.Error()}, fiber.StatusInternalServerError
	}
	authout := dto.AuthOut{Name: res.Name, DistribID: res.DistribID, AuthToken: tokenstring}
	return fiber.Map{"data": authout}, http.StatusAccepted

}
func GetUserByDistId(dist_id string) (fiber.Map, int) {

	var user models.User
	var result *gorm.DB

	user, result = repositories.GetUserByID(dist_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": user}, http.StatusOK
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
	nuser.BV = tc.Bv
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
	ruser.BV = tc.Bv

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

func FindNextAvailSlot(distrib_id string, center_code string, side string) (string, string) {
	/**
		On a Pyramid network, a reference can only be added either on left or right
		if a person adds thrid person an so on, the actual referree becomes the person below
		the person, if he has an empty slot on the same side
		if not the tree traverse till the botton where it finds an empty slot
		Here we find an empty slot recursively on the same side
		Caution a circular refernce by external db edit may cause an infinite loop
	**/
	fmt.Println("Finding Next Slot")
	var old_distrib_id string
	var old_center_code string
	for {
		old_distrib_id, old_center_code = distrib_id, center_code
		distrib_id, center_code = repositories.GetNextItem(distrib_id, center_code, side)
		fmt.Println("D: ", distrib_id, "C:", center_code)
		if distrib_id == "" {
			return old_distrib_id, old_center_code
		}
	}
}

func RegisterUser(user_in dto.UserIn) (fiber.Map, error) {
	//Generate Next Available Distrib Number
	distrib_id := FindNextAvailUserSeq()
	user := models.User{
		DistribID:     distrib_id,
		Name:          user_in.Name,
		Pass:          fmt.Sprintf("%x", sha256.Sum256([]byte(user_in.Pass))),
		Address1:      user_in.Address1,
		Address2:      user_in.Address2,
		TownOrCity:    user_in.TownOrCity,
		District:      user_in.District,
		EmailAddress:  user_in.EmailAddress,
		PinOrZipCode:  user_in.PinOrZipCode,
		Country:       user_in.Country,
		HomePhoneNo:   user_in.HomePhoneNo,
		MobilePhoneNo: user_in.MobilePhoneNo,
	}
	//if not empty don't overwrite but find next available free slot
	parent_distrib_id, parent_ref_place := FindNextAvailSlot(user_in.RefDistID, user_in.RefCenterCode, user_in.Place)

	tc1 := models.TrackingCenter{
		Name:           user_in.Name,
		DistribID:      distrib_id,
		CenterCode:     "001",
		PDistribId:     parent_distrib_id,
		PPlace:         parent_ref_place,
		LeftDistribID:  distrib_id,
		LeftPlace:      "002",
		RightDistribID: distrib_id,
		RightPlace:     "003",
	}
	tc2 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		CenterCode: "002",
		PDistribId: distrib_id,
		PPlace:     "001",
	}
	tc3 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		CenterCode: "003",
		PDistribId: distrib_id,
		PPlace:     "001",
	}

	//Succeed all or fail all
	tx := configs.DB.Begin()
	res := repositories.CreateUser(tx, user)
	if res != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error()}, res
	}
	res = repositories.CreateTCs(tx, []models.TrackingCenter{tc1, tc2, tc3})
	if res != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error()}, res
	}
	res = repositories.UpdateTC(tx, distrib_id, parent_distrib_id, parent_ref_place, user_in.Place)
	if res != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error()}, res
	}
	tx.Commit()
	rspdata := dto.UserOut{DistribID: distrib_id}
	return fiber.Map{"data": rspdata}, nil
}

func EditUserByDistId(DistribId string, userIn models.User) (fiber.Map, int) {

	var user models.User
	user, result := repositories.EditUserByDistId(DistribId, userIn, user)

	if result.Error != nil {
		log.Info("Error saving user to the database:", result.Error)
		return fiber.Map{"error": result.Error}, http.StatusBadGateway
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"success": "User Updated Successfully", "Deleted_User": user}, http.StatusOK
}

// Find Tree Up
func UpdateBvInTreeAfterPlaceOrder(orderId string, distrib_id string, place string, totalProductBV int) (fiber.Map, int) {

	var currentBV int

	for {

		//Getting Current BV from tc table
		currentBV, result := repositories.GetCurrentBvFromTc(distrib_id, place, currentBV)
		if result.Error != nil {
			return fiber.Map{"error": result.Error}, http.StatusInternalServerError
		}
		//adding currentbv and currently ordered product bv
		finalTotalBV := currentBV + totalProductBV

		//Updating each tc's BV from current tc to Upward tc
		if result := repositories.UpdateFinalBvToTc(distrib_id, place, finalTotalBV); result.Error != nil {
			return fiber.Map{"error": result.Error}, http.StatusInternalServerError
		}

		//record transaction in BV Transaction table
		tx := models.BvTransaction{
			DisribId:     distrib_id,
			Place:        place,
			OrderId:      orderId,
			Date:         time.Now(),
			BvValue:      finalTotalBV,
			ActivateDate: time.Now().AddDate(0, 0, 7),
		}
		//change name to add
		repositories.RecordBvTx(tx)

		//Get Next parent tracking center(upward)
		tc := repositories.GetTrackingCenter(distrib_id, place)
		//terminate condition of for loop
		if tc.PDistribId == "" {
			return fiber.Map{"data": "Success"}, http.StatusOK
		}
		//update parameters
		distrib_id = tc.PDistribId
		place = tc.PPlace
	}
}

func GetNewReferrals(distrib_id string) (fiber.Map, int) {

	var user []models.User
	var result *gorm.DB

	user, result = repositories.GetUserByRefDistribId(distrib_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": user}, http.StatusOK
}
