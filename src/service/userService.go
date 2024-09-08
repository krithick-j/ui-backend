package service

/**

A center code is a tc code which will have three values of 001,002,003 and will be created
by default. The same center code if referred in other places called place

**/

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func LoginUser(username string, password string, tx *gorm.DB) (fiber.Map, int) {
	pass := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	if username == "admin" {
		username = "IN-00001" //temporarily set IN-00001 as admin
	}

	res, err := repositories.AuthUser(username, pass, tx)
	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	if err != nil {
		configs.Log.Errorln("Error on calling AuthUser repositories fn from LoginUser service fn", err.Error())
		return fiber.Map{"err": err.Error()}, fiber.StatusUnauthorized
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

	//update last login
	err = repositories.UpdateLastLogin(res.DistribID, tx)
	if err != nil {
		utils.NotNilErrorMessage(err, "UpdateLastLogin", "LoginUser", fiber.StatusInternalServerError, tx)
	}

	authout := dto.AuthOut{
		Name:      res.Name,
		DistribID: res.DistribID,
		AuthToken: tokenstring,
		KYCStatus: res.KYCStatus,
		LastLogin: res.LastLogin.UTC(),
	}
	return fiber.Map{"data": authout}, fiber.StatusAccepted
}

func GetUserByDistribId(distribId string, tx *gorm.DB) (fiber.Map, int) {

	user, err := repositories.GetUserByID(distribId, tx)

	if user.DistribID == "" {
		configs.Log.Infoln("RecordNotFound")
		return fiber.Map{"data": "No user Found"}, fiber.StatusNoContent
	}

	if err != nil {
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GetUserByDistId fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": user}, fiber.StatusOK
}

func GetUsers(tx *gorm.DB) (fiber.Map, int) {
	users, err := repositories.GetAllUsers(tx)
	if err == gorm.ErrRecordNotFound {
		configs.Log.Infoln("RecordNotFound")
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		configs.Log.Errorln("Error on calling GetAllUsers repositories fn from GetUsers fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": users}, fiber.StatusOK
}

func GetUserRank(distribId string, tx *gorm.DB) (fiber.Map, int) {
	currentRank, titleRank, err := repositories.GetCurrentTitleRankByDistribId(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetCurrentTitleRankByDistribId", "GetUserRank", fiber.StatusInternalServerError, tx)
	}
	data := map[string]float64{
		"current_rank": currentRank,
		"title_rank":   titleRank,
	}
	return utils.SuccessMessage(data, fiber.StatusOK)
}

func FindNextAvailUserSeq(tx *gorm.DB) (string, error) {
	last_no, err := repositories.GetLastId(tx)
	if err != nil {
		configs.Log.Errorln("Error on calling GetLastId repositories fn from FindNextAvailUserSeq fn", err.Error())
		return err.Error(), err
	}
	distrib_no, _ := strconv.Atoi(strings.TrimPrefix(last_no, "IN-"))
	return fmt.Sprintf("IN-%05d", distrib_no+1), nil
}

func FindNextAvailSlot(distrib_id string, place string, side string, tx *gorm.DB) (string, string, error) {
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
		nextItem, err := repositories.GetNextItem(distrib_id, place, side, tx)
		if side == "left" {
			distrib_id = nextItem.LeftDistribID
			place = nextItem.LeftPlace
		} else if side == "right" {
			distrib_id = nextItem.RightDistribID
			place = nextItem.RightPlace
		} else {
			tx.Rollback()
			configs.Log.Errorln("Invalid Side output")
			return "", "", err
		}
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling CreateTCs repositories fn from Register User service", err.Error())
			return "", "", err
		}
		if distrib_id == "" {
			return old_distrib_id, old_place, nil
		}
	}
}

func RegisterUser(user_in dto.UserIn, tx *gorm.DB) (fiber.Map, int) {

	//Generate Next Available Distrib Number
	distrib_id, err := FindNextAvailUserSeq(tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling FindNextAvailUserSeq repositories fn from Register User service", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//user object
	user := handleUserRegistrationObject(user_in, distrib_id)

	//create tc and handles tc
	err = handleTCRegistration(user_in, distrib_id, user, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling handleTCRegistration repositories fn from Register User service", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//create distributor Application Form Link
	// res, err = CreateDistribApplicationForm(user)
	// if err != nil {
	// 	tx.Rollback()
	// 	configs.Log.Errorln("Error on calling CreateDistribApplicationForm repositories fn from Register User service", err.Error())
	// 	return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	// }

	//sends plain mail to user
	msg := fmt.Sprintf(`Dear Distributor, your registration in UI Network is successful. Your Distributor No is %s.`, distrib_id)
	SendPlainMail(user_in.EmailAddress, "Your Registration Details", msg)

	rspdata := dto.UserOut{DistribID: distrib_id}

	return fiber.Map{"data": rspdata}, fiber.StatusCreated
}

// func CreateDistribApplicationForm(user models.User, tx *gorm.DB) (fiber.Map, error) {
// 	//parse html template with values
// 	//convert html to pdf
// 	//save path in DistribApplicationFormLink
// 	return fiber.Map{"data": ""}, nil
// }

func handleUserRegistrationObject(user_in dto.UserIn, distrib_id string) models.User {
	AddressDetails := models.AddressDetails{
		Address1:        user_in.Address1,
		Address2:        user_in.Address2,
		TownOrCity:      user_in.TownOrCity,
		District:        user_in.District,
		StateOrProvince: user_in.StateOrProvince,
		PinOrZipCode:    user_in.PinOrZipCode,
		Country:         user_in.Country,
	}

	BankDetails := models.BankDetails{
		PanCard:   user_in.PanCard,
		BankName:  user_in.BankName,
		BankAccNo: user_in.BankAccNo,
		IFSCCode:  user_in.IFSCCode,
	}

	ApplicationInfo := models.ApplicationInformation{
		Title:                   user_in.Title,
		Name:                    user_in.Name,
		ChequeName:              user_in.ChequeName,
		EmailAddress:            user_in.EmailAddress,
		HomePhoneNo:             user_in.HomePhoneNo,
		MobilePhoneNo:           user_in.MobilePhoneNo,
		ValidIdNo:               user_in.ValidIdNo,
		DateOfBirth:             user_in.DateOfBirth,
		MothersMaidenName:       user_in.MothersMaidenName,
		BenificiaryName:         user_in.BenificiaryName,
		BeneficiaryRelationship: user_in.BeneficiaryRelationship,
		AddressDetails:          AddressDetails,
	}

	PreferredPlacementInformation := models.PreferredPlacementInformation{
		PreferredDistribId:   user_in.RefPlacementDistribId,
		PreferredDistribName: user_in.RefPlacementDistribname,
		PreferredPlace:       user_in.RefPlacementPlace,
		PreferredSide:        user_in.Side,
	}

	ReferrerInformation := models.ReferrerInformation{
		RefDistribID:   user_in.RefDistribID,
		RefDistribName: user_in.RefDistribName,
	}

	user := models.User{
		DistribID:                     distrib_id,
		Pass:                          fmt.Sprintf("%x", sha256.Sum256([]byte(user_in.Pass))),
		CpaPin:                        fmt.Sprintf("%x", sha256.Sum256([]byte(user_in.Pass))),
		ReferrerInformation:           ReferrerInformation,
		ApplicationInformation:        ApplicationInfo,
		BankDetails:                   BankDetails,
		PreferredPlacementInformation: PreferredPlacementInformation,
	}
	return user
}

func handleTCRegistration(user_in dto.UserIn, distrib_id string, user models.User, tx *gorm.DB) error {

	//if not empty don't overwrite but find next available free slot
	parent_distrib_id, parent_ref_place, err := FindNextAvailSlot(user_in.RefPlacementDistribId, user_in.RefPlacementPlace, user_in.Side, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling CreateUser repositories fn from Register User service", err.Error())
		return err
	}

	//Initial values
	place := "001"
	leftPlace := "002"
	rightPlace := "003"
	tc1 := models.TrackingCenter{
		Name:           user_in.Name,
		DistribID:      distrib_id,
		Place:          place,
		PDistribId:     parent_distrib_id,
		PPlace:         parent_ref_place,
		LeftDistribID:  distrib_id,
		LeftPlace:      leftPlace,
		RightDistribID: distrib_id,
		RightPlace:     rightPlace,
		IsActive:       false,
	}
	tc2 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      leftPlace,
		PDistribId: distrib_id,
		PPlace:     place,
		IsActive:   false,
	}
	tc3 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      rightPlace,
		PDistribId: distrib_id,
		PPlace:     place,
		IsActive:   false,
	}

	//Succeed all or fail all
	err = repositories.CreateUser(user, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling CreateUser repositories fn from Register User service", err.Error())
		return err
	}
	for _, newTc := range []models.TrackingCenter{tc1, tc2, tc3} {
		err = repositories.CreateTCs(newTc, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling CreateTCs repositories fn from Register User service", err.Error())
			return err
		}
	}

	//Updating Parent Tc after creating new TC
	err = repositories.UpdateTC(tx, distrib_id, place, parent_distrib_id, parent_ref_place, user_in.Side)
	if err != nil {
		configs.Log.Errorln("Error on calling UpdateTC repositories fn from Register User service", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func EditUserByDistId(DistribId string, userIn models.User, tx *gorm.DB) (fiber.Map, int) {
	if userIn.Pass != "" {
		return utils.CommonMessage("Cannot change password!", fiber.StatusForbidden, tx)
	}

	user, err := repositories.EditUserByDistId(DistribId, userIn, tx)
	if err != nil {
		configs.Log.Errorln("Error calling EditUserByDistId fn from EditUserByDistId service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"success": "User Updated Successfully", "user": user}, fiber.StatusOK
}

func GetNewReferrals(distrib_id string, tx *gorm.DB) (fiber.Map, int) {

	user, err := repositories.GetUserByRefDistribId(distrib_id, tx)

	if err == gorm.ErrRecordNotFound {
		configs.Log.Errorln("Error calling GetUserByRefDistribId fn from GetNewReferrals service fn", err.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		configs.Log.Errorln("Error calling GetUserByRefDistribId fn from GetNewReferrals service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": user}, fiber.StatusOK
}

func UpdateUserPass(payload dto.UserPassIn, tx *gorm.DB) (fiber.Map, int) {

	oldHashpassFromDB, err := repositories.GetUserPassByDistribId(payload.DistribId, tx)
	if err != nil {
		return fiber.Map{"err": err.Error()}, fiber.StatusInternalServerError
	}

	oldHashpass := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.OldPass)))
	if oldHashpass != oldHashpassFromDB {
		return fiber.Map{"data": "Old password does not match"}, fiber.StatusForbidden
	}

	newHashpass := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.NewPass)))

	err = repositories.UpdatePassword(payload.DistribId, newHashpass, tx)
	if err != nil {
		return fiber.Map{"err": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Password changed Successfully"}, fiber.StatusOK
}

func GetReferralChainByDistribId(distribId string, tx *gorm.DB) (fiber.Map, int) {
	var userArr []string

	// Recursive function to build the referral chain
	var buildReferralChain func(distribId string, tx *gorm.DB) error
	buildReferralChain = func(distribId string, tx *gorm.DB) error {
		referredDistribIds, err := repositories.GetReferredUsersByDistribId(distribId, tx)
		if err == gorm.ErrRecordNotFound {
			return nil // No more referrals found, end the recursion
		}
		if err != nil {
			return err // Return the error if something went wrong
		}

		for _, id := range referredDistribIds {
			userArr = append(userArr, id)    // Add to the chain
			err = buildReferralChain(id, tx) // Recurse to find the next level of referrals
			if err != nil {
				return err // Propagate the error if something went wrong during recursion
			}
		}
		return nil
	}

	// Start building the referral chain
	err := buildReferralChain(distribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetReferredUsersByDistribId", "GetReferralChainByDistribId", fiber.StatusInternalServerError, tx)
	}

	// Return the final referral chain
	return utils.SuccessMessage(userArr, fiber.StatusOK)
}
