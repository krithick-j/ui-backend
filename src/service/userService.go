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
	LeftPoint      int            `json:"left_point"`
	RightPoint     int            `json:"right_point"`
	BV             int            `json:"bv"`
	Left           *RecursiveUser `json:"left"`
	Right          *RecursiveUser `json:"right"`
}

func LoginUser(username string, password string) (fiber.Map, int) {
	pass := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if username == "admin" {
		username = "IN-00001" //temporarily set IN-00001 as admin
	}
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

func FindRecursiveTC(ruser *RecursiveUser, dist_id string, place string, side string) {
	tc := repositories.GetTrackingCenter(dist_id, place)
	var nuser *RecursiveUser = new(RecursiveUser)
	nuser.Name = tc.Name
	nuser.TrackingCenter = tc.DistribID + " " + tc.Place
	nuser.LeftPoint = tc.LeftPoint
	nuser.RightPoint = tc.RightPoint
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

func GetTreeUserByDistId(distrib_id string) fiber.Map {

	ruser := new(RecursiveUser)
	tc := repositories.GetTrackingCenter(distrib_id, "001")
	ruser.Name = tc.Name
	ruser.TrackingCenter = tc.DistribID + " " + tc.Place
	ruser.LeftPoint = tc.LeftPoint
	ruser.RightPoint = tc.RightPoint
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

func FindNextAvailSlot(distrib_id string, place string, side string) (string, string) {
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
	var old_place string
	for {
		old_distrib_id, old_place = distrib_id, place
		distrib_id, place = repositories.GetNextItem(distrib_id, place, side)
		fmt.Println("D: ", distrib_id, "C:", place)
		if distrib_id == "" {
			fmt.Printf("**************************************distrib %v place %v", old_distrib_id, old_place)
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
	}
	tc2 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      "002",
		PDistribId: distrib_id,
		PPlace:     "001",
	}
	tc3 := models.TrackingCenter{
		Name:       user_in.Name,
		DistribID:  distrib_id,
		Place:      "003",
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
	res = repositories.UpdateTC(tx, distrib_id, parent_distrib_id, parent_ref_place, user_in.Side)
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
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"success": "User Updated Successfully", "Deleted_User": user}, http.StatusOK
}

func UpdateCurrentPlaceValues(distrib_id string, placeBvs []dto.PlaceBv, orderId string) (fiber.Map, int) {
	var (
		tx = configs.DB.Begin()
		value int
		res *gorm.DB
	)
	//Adding Bv Points from the product to the tree
	for _, placeBv := range placeBvs {
		if placeBv.Place == "001" {
			//_, value = repositories.UpdateParentPlaceBv(distrib_id, placeBv)
		} else if placeBv.Place == "002" {
			res, value = repositories.UpdateLeftPointPlaceBv(distrib_id, placeBv, tx)
		} else if placeBv.Place == "003" {
			res, value = repositories.UpdateRightPointPlaceBv(distrib_id, placeBv, tx)
		} else {
			tx.Rollback()
			fmt.Println("Error in updating values of Place")
			return fiber.Map{"error": "Error in updating values of Place"}, http.StatusInternalServerError
		}
		if res.Error != nil {
			tx.Rollback()
			return fiber.Map{"error": "Error in updating values of Place"}, http.StatusInternalServerError
		}
		//record transaction in BV Transaction table
		BVtx := models.BvTransaction{
			DisribId:     distrib_id,
			Place:        placeBv.Place,
			OrderId:      orderId,
			Date:         time.Now(),
			BvValue:      value,
			ActivateDate: time.Now().AddDate(0, 0, 7),
		}
		//change name to add
		repositories.SaveBvTransaction(BVtx)
	}
	tx.Commit()
	return fiber.Map{"data": "Data Successfully Updated"}, http.StatusOK
}

func UpdateTreePlaceValuesByDistribId(distrib_id string) (fiber.Map, int) {

	var (
		place string = "001"
		tx = configs.DB.Begin()
	)

	for {
		//Get Next parent tracking center(upward)
		parentPlace := repositories.GetTrackingCenter(distrib_id, place)
		//terminate condition of for loop
		if parentPlace.PDistribId == "" {
			return fiber.Map{"data": "Data successfully updated in the tree"}, http.StatusOK
		}

		LeftDistribID := parentPlace.LeftDistribID
		leftPlace := parentPlace.LeftPlace
		RightDistribID := parentPlace.RightDistribID
		rightPlace := parentPlace.RightPlace

		//Get LeftTc
		leftTc := repositories.GetTrackingCenter(LeftDistribID, leftPlace)
		//Get RightTc
		rightTc := repositories.GetTrackingCenter(RightDistribID, rightPlace)
		//ParentPlace left and right point final values
		parentPlace.LeftPoint =  leftTc.RightPoint + leftTc.LeftPoint + leftTc.Bv
		parentPlace.RightPoint = rightTc.RightPoint + rightTc.LeftPoint + rightTc.Bv
		repositories.UpdateTrackingCenter(parentPlace.LeftPoint, parentPlace.RightPoint, parentPlace.DistribID, parentPlace.Place, tx)

		//update parameters
		distrib_id = parentPlace.PDistribId
		place = parentPlace.PPlace
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
