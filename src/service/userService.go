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
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func LoginUser(username string, password string) (fiber.Map, int) {
	pass := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if username == "admin" {
		username = "IN-00001" //temporarily set IN-00001 as admin
	}
	res, err := repositories.AuthUser(username, pass)
	if err != nil {
		configs.Log.Errorln("Error on calling AuthUser repositories fn from LoginUser service fn", err.Error())
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
		configs.Log.Errorln("Error on SignedString", err.Error())
		return fiber.Map{"err": err.Error()}, fiber.StatusInternalServerError
	}
	authout := dto.AuthOut{Name: res.Name, DistribID: res.DistribID, AuthToken: tokenstring, KYCStatus: res.KYCStatus}
	return fiber.Map{"data": authout}, http.StatusAccepted

}

func GetUserByDistId(dist_id string) (fiber.Map, int) {

	var user models.User
	var result *gorm.DB

	user, result = repositories.GetUserByID(dist_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("RecordNotFound")
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GetUserByDistId fn", result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": user}, fiber.StatusOK
}

func GetUsers() (fiber.Map, int) {
	var users []models.User
	var result *gorm.DB
	users, result = repositories.GetAllUsers(users)
	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("RecordNotFound")
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetAllUsers repositories fn from GetUsers fn", result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": users}, fiber.StatusOK
}

func FindNextAvailUserSeq() string {
	last_no := repositories.GetLastId()
	distrib_no, _ := strconv.Atoi(strings.TrimPrefix(last_no, "IN-"))
	return fmt.Sprintf("IN-%05d", distrib_no+1)
}

func FindNextAvailSlot(distrib_id string, place string, side string) (string, string) {
	/**
		On a Pyramid network, a reference can only be added either on left or right
		if a person adds thrid person and so on, the actual referree becomes the person below
		the person, if he has an empty slot on the same side
		if not the tree traverse till the bottoM where it finds an empty slot
		Here we find an empty slot recursively on the same side
		Caution a circular refernce by external db edit may cause an infinite loop
	**/
	var old_distrib_id string
	var old_place string
	for {
		old_distrib_id, old_place = distrib_id, place
		distrib_id, place = repositories.GetNextItem(distrib_id, place, side)
		if distrib_id == "" {
			return old_distrib_id, old_place
		}
	}
}

func RegisterUser(user_in dto.UserIn) (fiber.Map, error) {
	//Generate Next Available Distrib Number
	distrib_id := FindNextAvailUserSeq()
	user := models.User{
		Name:            user_in.Name,
		Pass:            fmt.Sprintf("%x", sha256.Sum256([]byte(user_in.Pass))),
		RefDistribID:    user_in.RefDistribID,
		DistribID:       distrib_id,
		Address1:        user_in.Address1,
		Address2:        user_in.Address2,
		TownOrCity:      user_in.TownOrCity,
		District:        user_in.District,
		StateOrProvince: user_in.StateOrProvince,
		EmailAddress:    user_in.EmailAddress,
		PinOrZipCode:    user_in.PinOrZipCode,
		Country:         user_in.Country,
		HomePhoneNo:     user_in.HomePhoneNo,
		MobilePhoneNo:   user_in.MobilePhoneNo,
	}
	//if not empty don't overwrite but find next available free slot
	parent_distrib_id, parent_ref_place := FindNextAvailSlot(user_in.RefPlacementDistribId, user_in.RefPlacementPlace, user_in.Side)

	tc1 := models.TrackingCenter{
		Name:           user_in.Name,
		DistribID:      distrib_id,
		Place:          "001",
		PDistribId:     parent_distrib_id,
		PPlace:         parent_ref_place,
		LeftDistribID:  distrib_id,
		LeftPlace:      "002",
		RightDistribID: distrib_id,
		RightPlace:     "003",
		IsActive:       false,
	}
	tc2 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      "002",
		PDistribId: distrib_id,
		PPlace:     "001",
		IsActive:   false,
	}
	tc3 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      "003",
		PDistribId: distrib_id,
		PPlace:     "001",
		IsActive:   false,
	}

	//Succeed all or fail all
	tx := configs.DB.Begin()
	res := repositories.CreateUser(tx, user)
	if res != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling CreateUser repositories fn from Register User service", res.Error())
		return fiber.Map{"error": res.Error()}, res
	}
	res = repositories.CreateTCs(tx, []models.TrackingCenter{tc1, tc2, tc3})
	if res != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling CreateTCs repositories fn from Register User service", res.Error())
		return fiber.Map{"error": res.Error()}, res
	}
	res = repositories.UpdateTC(tx, distrib_id, parent_distrib_id, parent_ref_place, user_in.Side)
	if res != nil {
		configs.Log.Errorln("Error on calling UpdateTC repositories fn from Register User service", res.Error())
		tx.Rollback()
		return fiber.Map{"error": res.Error()}, res
	}
	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on Committing Transaction RegisterUser service fn", err.Error())
		tx.Rollback()
		return fiber.Map{"Error": err.Error()}, err

	}
	rspdata := dto.UserOut{DistribID: distrib_id}
	msg := fmt.Sprintf(`Dear User,
We received your signup. Your login id/distributor id is %s.
Please login and upload your KYC documents to proceed.`, distrib_id)
	SendPlainMail(user_in.EmailAddress, "Your Signup Details", msg)
	return fiber.Map{"data": rspdata}, nil
}

func EditUserByDistId(DistribId string, userIn models.User) (fiber.Map, int) {

	var user models.User
	user, result := repositories.EditUserByDistId(DistribId, userIn, user)

	if result.Error != nil {
		configs.Log.Errorln("Error calling EditUserByDistId fn from EditUserByDistId service fn", result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}
	return fiber.Map{"success": "User Updated Successfully", "Deleted_User": user}, fiber.StatusOK
}

func GetNewReferrals(distrib_id string) (fiber.Map, int) {

	var user []models.User
	var result *gorm.DB

	user, result = repositories.GetUserByRefDistribId(distrib_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("Error calling GetUserByRefDistribId fn from GetNewReferrals service fn", result.Error.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error calling GetUserByRefDistribId fn from GetNewReferrals service fn", result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": user}, fiber.StatusOK
}

func UpdateUserPass(payload dto.UserPassIn) (fiber.Map, int) {

	oldHashpassFromDB, res := repositories.GetUserPassByDistribId(payload.DistribId)
	if res.Error != nil {
		return fiber.Map{"err": res.Error.Error()}, fiber.StatusInternalServerError
	}

	oldHashpass := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.OldPass)))
	if oldHashpass != oldHashpassFromDB {
		return fiber.Map{"data": "Old password does not match"}, fiber.StatusForbidden
	}

	newHashpass := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.NewPass)))

	res = repositories.UpdatePassword(payload.DistribId, newHashpass)
	if res.Error != nil {
		return fiber.Map{"err": res.Error.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": "Password changed Successfully"}, http.StatusOK

}
