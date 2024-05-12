package service

/**

A center code is a tc code which will have three values of 001,002,003 and will be created
by default. The same center code if referred in other places called place

**/

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"path/filepath"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	authout := dto.AuthOut{Name: res.Name, DistribID: res.DistribID, AuthToken: tokenstring, KYCStatus: res.KYCStatus}
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
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fiber.Map{"Error": err.Error()}, err

	}
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

func GenerateIDCard(distrib_id string) error {
	configs.Log.Info("Started Generating ID Card")
	user := models.User{}
	userdata, _ := repositories.GetUserByID(distrib_id, user)
	userphoto, _ := strings.CutPrefix(userdata.KYCPhoto, "media/")
	extension := filepath.Ext(userphoto)
	extension, _ = strings.CutPrefix(extension, ".")
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// Outer Rect
	pdf.Rect(5, 5, 180, 60, "D")
	// Inner Left Rect
	pdf.Rect(10, 10, 80, 50, "D")
	// Inner Right Rect
	pdf.Rect(100, 10, 80, 50, "D")
	//Disclaimer Box
	pdf.SetFillColor(234, 234, 245)
	pdf.RoundedRect(48, 38, 40, 15, 1, "1234", "DF")
	//Address Footer
	pdf.Rect(100, 45, 80, 15, "DF")

	userphoto = filepath.Join("./assets", userphoto)
	if _, err := os.Stat(userphoto); errors.Is(err, os.ErrNotExist) {
		configs.Log.Errorf("File Does Not exit %s", userphoto)
	}
	if _, err := os.Stat("./assets/images/uilogo.png"); errors.Is(err, os.ErrNotExist) {
		configs.Log.Error("File Does Not exit", zap.String("image", "./assets/images/uilogo.png"))
	}
	pdf.Image(userphoto, 15, 10, 25, 0, false, extension, 0, "")
	pdf.Image("./assets/images/uilogo.png", 70, 12, 6, 0, false, "png", 0, "")
	pdf.Image("./assets/images/uilogo.png", 104, 49, 6, 0, false, "png", 0, "")
	// Inner Right Rect
	pdf.SetFillColor(0, 255, 0)
	pdf.RoundedRect(11, 42, 32, 6, 2, "1234", "DF")
	pdf.SetFont("Arial", "B", 12)
	var (
		currX float64 = 12
		currY float64 = 45
	)

	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "UI Distributor")
	currX += 8
	currY += 7
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Of")
	currX += -8
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Universe International - India")
	pdf.SetFont("Arial", "B", 10)
	currX += 36
	currY -= 34
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, userdata.Name)
	pdf.SetFont("Arial", "", 8)
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Dist ID: "+userdata.DistribID)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Email: "+userdata.EmailAddress)
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Phone: +91-"+userdata.MobilePhoneNo)
	pdf.SetFont("Arial", "I", 8)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Registration Fees"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Deposits"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Investments"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• No Job Offerings"))
	/*
		Back of the card
	*/
	currX += 55
	currY -= 34
	pdf.SetFont("Arial", "", 6)
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "Phone No: 9080296128")
	currY += 4
	//pdf.SetXY(currX, currY)
	//pdf.Cell(0, 0, "CIN: XXXXXXX-XXXXX")
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "GST: Applied")
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• Check Government ID Card for Proof"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• This card is not transferable"))
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, tr("• This card should be shown before presentation"))
	currY += 5
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(0, 0, tr("In Case of Any Enquiry, Write to:"))
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "support@ui-network.com")
	//Footer
	currX += 10
	currY += 6
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 0, "Universe International - India")
	currY += 4
	pdf.SetXY(currX, currY)
	pdf.SetFont("Arial", "", 6)
	pdf.Cell(0, 0, "KASI ARCADE FIRST FLOOR, VOC STREET,KAIKANKUPPAM")
	currY += 3
	pdf.SetXY(currX, currY)
	pdf.Cell(0, 0, "RAMAPURAM, Chennai - PIN 600087 - TAMIL NADU")
	filename := fmt.Sprintf("./assets/%s-idcard.pdf", distrib_id)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		configs.Log.Error("Error is ", zap.String("outfile", err.Error()))
		return err
	}
	return nil
}

func SendEmailCode(toMail string) error {
	otp := GenOPT()
	err := repositories.SaveOTP("email", toMail, otp)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Your OTP for Email Verification is %s", otp)
	tomail := fmt.Sprintf("<%s>", toMail)
	SendMail(tomail, "Your Email Verification OTP", mail.TypeTextPlain, msg)
	return nil
}

func VerifyEmailCode(otp string, email string) error {
	return repositories.CheckAndUpdateOTP("email", email, otp)
}

func SendPhoneCode(c *fiber.Ctx, phone string) error {
	otp := GenOPT()
	url := "https://www.textguru.in/api/v22.0/?"
	payload := fmt.Sprintf("username=jega.in&password=60423479&source=GSENTS&dmobile=91%s&dlttempid=1707171500974884924&message=Dear Customer,\nThis is your OTP for Login %s for your mobile number verification On https://ui-network.com.\nGSENTS", phone, otp)
	agent := fiber.Post(url)
	agent.Body([]byte(payload)) // set body received by request
	statusCode, body, errs := agent.Bytes()
	_, _ = statusCode, body
	defer agent.ConnectionClose()
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("Error is ", err.Error())
		}
	}
	err := repositories.SaveOTP("phone", phone, otp)
	if err != nil {
		return err
	}
	return nil
}

func VerifyPhoneCode(otp string, phone string) error {
	return repositories.CheckAndUpdateOTP("phone", phone, otp)
}

func KycUpload(c *fiber.Ctx, form *multipart.Form) error {
	user := models.User{}
	user.DistribID = form.Value["distrib_id"][0]
	for fs, fhs := range form.File {
		for _, fh := range fhs {
			extension := filepath.Ext(fh.Filename)
			fullPath := "./assets/" + user.DistribID + "-" + fs + extension
			mediapath := "media/" + user.DistribID + "-" + fs + extension
			err := c.SaveFile(fh, fullPath)
			if err != nil {
				return err
			}
			switch fs {
			case "aadhar":
				user.KYCAdhaar = mediapath
			case "consent":
				user.KYCConsentDoc = mediapath
			case "pan":
				user.KYCPAN = mediapath
			case "user-image":
				user.KYCPhoto = mediapath
			}
		}
	}
	user.KYCStatus = "pending"
	repositories.UpdateKyc(&user)
	return nil
}

func ApproveKYC(distrib_id string) error {
	user := models.User{DistribID: distrib_id, KYCStatus: "verified"}
	repositories.UpdateKyc(&user)
	return nil
}
