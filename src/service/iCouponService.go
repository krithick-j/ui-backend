package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Generates Unique Hex Code of length 10. Length != 10 return error
func GenerateUniqueHexCode(length int) string {
	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	hexCode := hex.EncodeToString(randomBytes)
	hexCode = hexCode[:length]

	return strings.ToUpper(hexCode)
}

//func SendICouponMail(toMail string, coupons []dto.SendCoupon) error {

// Parse the email template
// tmpl, err := template.ParseFiles("assets/templates/email_template.html")
// if err != nil {
// 	return err
// }

// // Create a new file to store the rendered email content
// file, err := os.Create("email.html")
// if err != nil {
// 	return err
// }
// defer file.Close()

// Execute the template with the coupons data and write it to the file
// err = tmpl.Execute(file, coupons)
// if err != nil {
// 	return err
// }

// Read the contents of the rendered HTML file
// renderedEmail, err := os.ReadFile("email.html")
// if err != nil {
// 	return err
// }

// m, err := gomail.NewGoMail()
// if err != nil {
// 	return err
// }

// m.Set("Username", "j.krithick@gmail.com")
// m.Set("Password", "gior zeiv xgga lqky")

// m.Set("Servername", "smtp.gmail.com:465")

// m.Set("From", "j.krithick@gmail.com")
// m.Set("From_name", "Krithick ")

// m.Set("To", toMail)

// m.Set("Subject", "Your new iCoupon")

// m.Set("BodyMessage", string(renderedEmail))

// if err := m.SendMessage(); err != nil {
// 	return err
// }
// 	return nil
// }

// This function is used to generate ICoupon and send email
func AddICoupon(iCouponIn dto.ICouponIn, adminName string,expireDate time.Time, tx *gorm.DB) (fiber.Map, int) {

	var iCoupon models.ICoupon
	var iCoupons []dto.SendCoupon
	reference := GenerateUniqueHexCode(10)

	email, err := repositories.GetUserEmailByDistribID(iCouponIn.DistribID, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling  GetUserEmailByDistribID repositories fn from AddICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	if email == "" {
		configs.Log.Errorln("RecordNotFound on calling  GetUserEmailByDistribID repositories fn from AddICoupon service fn")
		tx.Rollback()
		return fiber.Map{"data": "No emailID Found"}, fiber.StatusNotFound
	}

	for _, Coupon := range iCouponIn.Coupons {
		for i := 0; i < int(Coupon.Quantity); i++ {
			hexVID := GenerateUniqueHexCode(10)
			hexPin := GenerateUniqueHexCode(10)
			iCoupon = models.ICoupon{
				VID:       hexVID,
				Pin:       hexPin,
				DistribID: iCouponIn.DistribID,
				DateOn:    time.Now(),
				Reference: reference,
				ExpiresOn: expireDate,
				Value:     Coupon.Value,
				AdminName: adminName,
				Active:    true,
			}

			SendCoupon := dto.SendCoupon{
				VID:       iCoupon.VID,
				Pin:       iCoupon.Pin,
				Value:     iCoupon.Value,
				DateOn:    utils.FormatTimeByLocation(iCoupon.ExpiresOn, "Asia/Kolkata", "02-01-2006"),
				ExpiresOn: utils.FormatTimeByLocation(iCoupon.ExpiresOn, "Asia/Kolkata", "02-01-2006"),
				Active:    iCoupon.Active,
			}
			iCoupons = append(iCoupons, SendCoupon)
			err := repositories.SaveICoupon(iCoupon, tx)
			if err != nil {
				tx.Rollback()
				configs.Log.Errorln("Error on calling SaveICoupon repositories fn from AddICoupon service fn", err.Error())
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}

			ICouponTxObj := models.ICouponTransaction{
				DistribId: iCoupon.DistribID,
				VID:       iCoupon.VID,
				Value:     iCoupon.Value,
				Reference: reference,
			}
			err = repositories.SaveICouponTx(ICouponTxObj, tx)
			if err != nil {
				tx.Rollback()
				configs.Log.Errorln("Error on calling SaveICouponTx repositories fn from AddICoupon service fn", err.Error())
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}
		}
	}
	err = SendHtmlMailICoupon(email, "Your new iCoupon", iCoupons)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SendHtmlMailICoupon service fn from AddICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "ICoupons added successfully and sent to your mail"}, fiber.StatusCreated
}

func GetAllICouponsByDistribId(DistribID string, tx *gorm.DB) (fiber.Map, int) {
	allICouponsOut := make([]dto.GetICouponOut, 0)
	iCoupons, err := repositories.GetAllICouponsByDistribID(DistribID, tx)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.
			Errorln("RecordNotFound on calling  GetAllICouponsByDistribID repositories fn from GetAllICouponsByDistribId service fn", err.Error())
		return fiber.Map{"data": "No ICoupon exists"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling  GetAllICouponsByDistribID repositories fn from GetAllICouponsByDistribId service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	for _, iCoupon := range iCoupons {

		iCouponBalance, err := repositories.GetICouponBalance(iCoupon.VID, tx)
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			configs.Log.Errorln("RecordNotFound on calling  ICouponBalance repositories fn from GetAllICouponsByDistribId service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusNotFound
		}
		if err != nil {
			tx.Rollback()
			configs.Log.
				Errorln("Error on calling  ICouponBalance repositories fn from GetAllICouponsByDistribId service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		iCouponDto := dto.GetICouponOut{
			DistribID:      DistribID,
			DateOn:         iCoupon.DateOn.UTC().String(),
			Reference:      iCoupon.Reference,
			AdminName:      iCoupon.AdminName,
			VID:            iCoupon.VID,
			TotalValue:     iCoupon.Value,
			RemainingValue: iCouponBalance,
			ExpiresOn:      iCoupon.ExpiresOn.UTC().String(),
			Pin:            iCoupon.Pin,
			Active:         iCoupon.Active,
			CreatedAt:      iCoupon.CreatedAt.UTC().String(),
			UpdatedAt:      iCoupon.UpdatedAt.UTC().String(),
			ID:             iCoupon.ID,
		}
		allICouponsOut = append(allICouponsOut, iCouponDto)
	}
	return fiber.Map{"data": allICouponsOut}, fiber.StatusOK
}

func ValidateICoupon(payload dto.ValidateICouponIn, distribID string, tx *gorm.DB) (fiber.Map, int) {
	iCoupon, err := repositories.ValidateICoupon(payload.VID, payload.Pin, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling  ValidateICoupon repositories fn from ValidateICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	if !iCoupon.Active {
		return fiber.Map{"data": "Icoupon expired"}, fiber.StatusBadRequest
	}
	//date expiry condition
	expiresOn, err := repositories.GetICouponExpiryDate(payload.VID, tx)
	if err != nil {
		configs.Log.Errorln("Error on calling GetICouponExpiryDate fn from ValidateICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusBadRequest
	}
	if expiresOn.Before(time.Now()) || expiresOn.Equal(time.Now()) {
		repositories.CloseCoupon(payload.VID, tx)
		configs.Log.Errorln("ICoupon Expired on ", expiresOn.String())
		return fiber.Map{"error": fmt.Sprintf("ICoupon expired on %s", expiresOn.String())}, fiber.StatusBadRequest
	}

	balance, err := repositories.GetICouponBalance(payload.VID, tx)

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.
			Errorln("RecordNotFound on calling  GetICouponBalance repositories fn from ValidateICoupon service fn", err.Error())
		return fiber.Map{"data": "No ICoupon exists"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetICouponBalance fn from ValidateICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusBadRequest
	}

	iCouponsOut := dto.ValidateICouponOut{
		Value:      balance, //Remaining Value
		TotalValue: iCoupon.Value,
	}

	return fiber.Map{"data": iCouponsOut}, fiber.StatusOK
}

func GetICouponHistory(payload dto.ICouponHistoryIn, tx *gorm.DB) (fiber.Map, int) {
	var fromDate, toDate time.Time

	if payload.FromDate == "" || payload.ToDate == "" {
		iCouponHistory, err := repositories.GetICouponHistory(payload.DistribId, tx)
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			configs.Log.
				Errorln("No ICoupon Transaction History on calling GetICouponHistory repositories fn from GetICouponHistory service fn ", err.Error())
			return fiber.Map{"data": "No ICoupon Transaction History"}, fiber.StatusNotFound
		}
		if err != nil {
			tx.Rollback()
			configs.Log.
				Errorln("Error on calling GetICouponHistory repositories fn from GetICouponHistory service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
		return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
	}

	fromDate, err := time.Parse("2006-01-02", payload.FromDate)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling Parse fn from GetICouponHistory service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	toDate, err = time.Parse("2006-01-02", payload.ToDate)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling Parse fn from GetICouponHistory service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	// Adjust time to start of the day (00:00:00) for fromDate and end of the day (23:59:59) for toDate
	fromDateStart := fromDate.Format("2006-01-02 15:04:05")
	toDateEnd := toDate.Add(24*time.Hour - time.Second).Format("2006-01-02 15:04:05")

	iCouponHistory, err := repositories.GetICouponHistoryByDate(payload.DistribId, fromDateStart, toDateEnd, tx)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Infoln("No record Found on calling GetICouponHistoryByDate repositories fn from GetICouponHistory service fn", err.Error())
		return fiber.Map{"data": "Record Not Found on calling GetICouponHistoryByDate repostories fn from GetICouponHistory service fn"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetICouponHistoryByDate repositories fn from GetICouponHistory service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
}

// This function is used to get total ICoupon and Total Used value
func GetICouponArrayByOrderId(orderId string, distribId string, tx *gorm.DB) (fiber.Map, int) {

	//Creating Output dto obj
	ObjOut := []dto.OrderedICouponOut{}
	//Getting ICoupon Number
	getICouponsVID, err := repositories.GetICouponsVIDByReference(orderId, distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SaveICoupon repositories fn from AddICoupon service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//Rotating all the VID and mapping total balance and remaining value
	for _, vid := range getICouponsVID {
		totalValue, err := repositories.GetICouponValueByVID(vid, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetICouponValueByVID repositories fn from GetICouponArrayByOrderId service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
		//get ICoupon balance
		balance, err := repositories.GetICouponBalance(vid, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetICouponBalance repositories fn from GetICouponArrayByOrderId service fn")
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		icoupon := dto.OrderedICouponOut{
			VID:       vid,
			Value:     totalValue,
			UsedValue: totalValue - balance,
		}

		ObjOut = append(ObjOut, icoupon)
		//Getting Total Balance By ICoupon Number
	}
	return fiber.Map{"data": ObjOut}, fiber.StatusOK
}

// This function is used to get total ICoupon and Total Used value
func GetICouponArrayTotalValueByOrderId(orderId string, distribId string, tx *gorm.DB) (fiber.Map, int) {

	//Creating Output dto obj
	ObjOut := []dto.OrderedICouponOut{}
	//Getting ICoupon Number
	getICouponsVID, err := repositories.GetICouponsVIDByReference(orderId, distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SaveICoupon repositories fn from GetICouponArrayTotalValueByOrderId service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//Rotating all the VID and mapping total balance and remaining value
	for _, vid := range getICouponsVID {
		totalValue, err := repositories.GetICouponValueByVID(vid, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetICouponValueByVID repositories fn from GetICouponArrayTotalValueByOrderId service fn", err.Error())
			return fiber.Map{"error": err.Error(), "err": err}, fiber.StatusInternalServerError
		}

		//get ICoupon balance By Order
		UsedValue, err := repositories.GetICouponRowsByOrderID(orderId, vid, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetICouponRowsByOrderID fn from GetICouponArrayTotalValueByOrderId service fn")
			return fiber.Map{"data": ObjOut}, fiber.StatusInternalServerError
		}

		icoupon := dto.OrderedICouponOut{
			VID:       vid,
			Value:     totalValue,
			UsedValue: math.Abs(UsedValue),
		}

		ObjOut = append(ObjOut, icoupon)
		//Getting Total Balance By ICoupon Number
	}
	return fiber.Map{"data": ObjOut}, fiber.StatusOK
}
