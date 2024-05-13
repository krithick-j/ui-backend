package service

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
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

	result, email := repositories.GetUserEmailByDistribID(iCouponIn.DistribID)

	if result.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fiber.Map{"data": "No emailID Found"}, http.StatusNotFound
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
				tx.Rollback()
				return fiber.Map{"error": err.Error()}, http.StatusInternalServerError
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
				return fiber.Map{"error": res.Error.Error()}, http.StatusInternalServerError
			}
		}
	}
	SendHtmlMail(email, "Your new iCoupon", iCoupons)

	return fiber.Map{"data": "ICoupons added successfully and sent to your mail"}, http.StatusCreated
}

func GetAllICouponsByDistribId(DistribID string) (fiber.Map, int) {
	var iCoupons []models.ICoupon
	result, iCoupons := repositories.GetAllICouponsByDistribID(DistribID, iCoupons)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon exists"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": iCoupons}, http.StatusOK

}

func ValidateICoupon(payload dto.ValidateICouponIn, distribID string) (fiber.Map, int) {
	var iCoupon models.ICoupon
	iCoupon = repositories.ValidateICoupon(payload.VID, payload.Pin, iCoupon)

	// if iCoupon.DistribID != distribID {
	// 	return fiber.Map{"data": "Invalid ICoupon"}, http.StatusForbidden
	// }

	if !iCoupon.Active {
		return fiber.Map{"data": "Icoupon expired"}, http.StatusOK
	}
	//date expiry condition
	balance, result := repositories.GetICouponBalance(payload.VID)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon exists"}, http.StatusNotFound
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	iCouponsOut := dto.ValidateICouponOut{
		Value: balance,
	}

	return fiber.Map{"data": iCouponsOut}, http.StatusOK
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
