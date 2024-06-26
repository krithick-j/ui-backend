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
func AddICoupon(iCouponIn dto.ICouponIn, adminName string, tx *gorm.DB) (fiber.Map, int) {

	var iCoupon models.ICoupon
	var iCoupons []dto.SendCoupon
	reference := GenerateUniqueHexCode(10)

	email, result := repositories.GetUserEmailByDistribID(iCouponIn.DistribID)

	if result.Error != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling  GetUserEmailByDistribID repositories fn from AddICoupon service fn", result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		configs.Log.Errorln("RecordNotFound on calling  GetUserEmailByDistribID repositories fn from AddICoupon service fn", result.Error.Error())
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
				ExpiresOn: time.Now().AddDate(0, 6, 0),
				Value:     Coupon.Value,
				AdminName: adminName,
				Active:    true,
			}

			SendCoupon := dto.SendCoupon{
				VID:       iCoupon.VID,
				Pin:       iCoupon.Pin,
				Value:     iCoupon.Value,
				DateOn:    iCoupon.DateOn,
				ExpiresOn: iCoupon.ExpiresOn,
				Active:    iCoupon.Active,
			}
			iCoupons = append(iCoupons, SendCoupon)
			err := repositories.SaveICoupon(iCoupon)
			if err != nil {
				configs.Log.Errorln("Error on calling SaveICoupon repositories fn from AddICoupon service fn", result.Error.Error())
				tx.Rollback()
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}

			ICouponTxObj := models.ICouponTransaction{
				DistribId: iCoupon.DistribID,
				VID:       iCoupon.VID,
				Value:     iCoupon.Value,
				Reference: reference,
			}
			res := repositories.SaveICouponTx(ICouponTxObj)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		}
	}
	fmt.Println("email--->", email)
	SendHtmlMailICouopon(email, "Your new iCoupon", iCoupons)

	return fiber.Map{"data": "ICoupons added successfully and sent to your mail"}, fiber.StatusCreated
}

func GetAllICouponsByDistribId(DistribID string) (fiber.Map, int) {
	allICouponsOut := make([]dto.GetICouponOut, 0)
	fmt.Println("all icoupons", allICouponsOut)
	result, iCoupons := repositories.GetAllICouponsByDistribID(DistribID)
	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("RecordNotFound on calling  GetAllICouponsByDistribID repositories fn from GetAllICouponsByDistribId service fn", result.Error.Error())
		return fiber.Map{"data": "No ICoupon exists"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling  GetAllICouponsByDistribID repositories fn from GetAllICouponsByDistribId service fn", result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}

	for _, iCoupon := range iCoupons {

		iCouponBalance, res := repositories.GetICouponBalance(iCoupon.VID)
		if res.Error == gorm.ErrRecordNotFound {
			configs.Log.Errorln("RecordNotFound on calling  ICouponBalance repositories fn from GetAllICouponsByDistribId service fn", res.Error.Error())
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusNotFound
		}

		if res.Error != nil {
			configs.Log.Errorln("Error on calling  ICouponBalance repositories fn from GetAllICouponsByDistribId service fn", res.Error.Error())
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
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

func ValidateICoupon(payload dto.ValidateICouponIn, distribID string) (fiber.Map, int) {
	var iCoupon models.ICoupon
	iCoupon = repositories.ValidateICoupon(payload.VID, payload.Pin, iCoupon)
	fmt.Println("icouponn------>>", iCoupon)
	if !iCoupon.Active {
		return fiber.Map{"data": "Icoupon expired"}, fiber.StatusBadRequest
	}
	//date expiry condition
	expiresOn, result := repositories.GetICouponExpiryDate(payload.VID)
	if result.Error != nil {
		configs.Log.Errorln("Error Retrieving the ICoupon Expiry Date")
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusBadRequest
	}

	if expiresOn.Before(time.Now()) || expiresOn.Equal(time.Now()) {
		repositories.CloseCoupon(payload.VID)
		configs.Log.Errorln("ICoupon Expired on ", expiresOn.String())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusBadRequest
	}

	balance, result := repositories.GetICouponBalance(payload.VID)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon exists"}, fiber.StatusNotFound
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}

	iCouponsOut := dto.ValidateICouponOut{
		Value:      balance, //Remaining Value
		TotalValue: iCoupon.Value,
	}

	fmt.Println("data", iCouponsOut)

	return fiber.Map{"data": iCouponsOut}, fiber.StatusOK
}

func GetICouponHistory(payload dto.ICouponHistoryIn) (fiber.Map, int) {
	var fromDate, toDate time.Time

	if payload.FromDate == "" || payload.ToDate == "" {
		iCouponHistory, err := repositories.GetICouponHistory(payload.DistribId)

		if err == gorm.ErrRecordNotFound {
			return fiber.Map{"data": "No ICoupon Transaction History"}, fiber.StatusNotFound
		}
		return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
	}

	fromDate, err := time.Parse("2006-01-02", payload.FromDate)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusNotFound
	}

	toDate, err = time.Parse("2006-01-02", payload.ToDate)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusNotFound
	}

	// Adjust time to start of the day (00:00:00) for fromDate and end of the day (23:59:59) for toDate
	fromDateStart := fromDate.Format("2006-01-02 15:04:05")
	toDateEnd := toDate.Add(24*time.Hour - time.Second).Format("2006-01-02 15:04:05")

	iCouponHistory, err := repositories.GetICouponHistoryByDate(payload.DistribId, fromDateStart, toDateEnd)

	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon Transaction History"}, fiber.StatusNotFound
	}
	return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
}

// This function is used to get total ICoupon and Total Used value
func GetICouponArrayByOrderId(orderId string, distribId string) ([]dto.OrderedICouponOut, error) {

	//Creating Output dto obj
	ObjOut := []dto.OrderedICouponOut{}
	//Getting ICoupon Number
	getICouponsVID, err := repositories.GetICouponsVIDByReference(orderId, distribId)
	if err != nil {
		return ObjOut, err
	}

	//Rotating all the VID and mapping total balance and remaining value
	for _, vid := range getICouponsVID {
		totalValue := repositories.GetICouponValueByVID(vid)

		//get ICoupon balance
		balance, result := repositories.GetICouponBalance(vid)
		if result.Error != nil {
			configs.Log.Errorln("Error Retrieving the ICoupon Balance")
			return ObjOut, err
		}

		icoupon := dto.OrderedICouponOut{
			VID:       vid,
			Value:     totalValue,
			UsedValue: totalValue - balance,
		}

		ObjOut = append(ObjOut, icoupon)
		//Getting Total Balance By ICoupon Number
	}
	return ObjOut, nil
}

// This function is used to get total ICoupon and Total Used value
func GetICouponArrayTotalValueByOrderId(orderId string, distribId string) ([]dto.OrderedICouponOut, error) {

	//Creating Output dto obj
	ObjOut := []dto.OrderedICouponOut{}
	//Getting ICoupon Number
	getICouponsVID, err := repositories.GetICouponsVIDByReference(orderId, distribId)
	if err != nil {
		return ObjOut, err
	}

	//Rotating all the VID and mapping total balance and remaining value
	for _, vid := range getICouponsVID {
		totalValue := repositories.GetICouponValueByVID(vid)

		//get ICoupon balance By Order
		UsedValue, result := repositories.GetICouponRowsByOrderID(orderId)
		if result.Error != nil {
			configs.Log.Errorln("Error Retrieving the ICoupon Balance")
			return ObjOut, err
		}

		icoupon := dto.OrderedICouponOut{
			VID:       vid,
			Value:     totalValue,
			UsedValue: math.Abs(UsedValue),
		}

		ObjOut = append(ObjOut, icoupon)
		//Getting Total Balance By ICoupon Number
	}
	return ObjOut, nil
}
