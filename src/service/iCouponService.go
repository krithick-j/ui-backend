package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/fmorenovr/gomail"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func generateUniqueHexCode(length int) string {
	randomBytes := make([]byte, length/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	hexCode := hex.EncodeToString(randomBytes)
	hexCode = hexCode[:length]

	return strings.ToUpper(hexCode)
}

func SendICouponMail(toMail string, coupons []dto.SendCoupon) error {
	m, err := gomail.NewGoMail()
	if err != nil {
		return err
	}

	m.Set("Username", "j.krithick@gmail.com")
	m.Set("Password", "gior zeiv xgga lqky")

	m.Set("Servername", "smtp.gmail.com:465")

	m.Set("From", "j.krithick@gmail.com")
	m.Set("From_name", "Krithick ")

	m.Set("To", toMail)

	m.Set("Subject", "Your new iCoupon")
	fmt.Println(coupons)
	// Construct the body message
	var body strings.Builder
	body.WriteString("This is a noreply email. Your iCoupons are:\n\n")
	body.WriteString("<tr><th>VID</th><th>PIN</th><th>Value</th><th>Date On</th><th>Expires On</th></tr>\n")
	for _, coupon := range coupons {
		DateOnFormat := coupon.DateOn.Format("02-01-06")
		ExpireOnFormat := coupon.ExpiresOn.Format("02-01-06")
		body.WriteString("<tr>")
		body.WriteString("<td>" + coupon.VID + "</td>")
		body.WriteString("<td>" + coupon.Pin + "</td>")
		body.WriteString("<td>" + strconv.FormatFloat(coupon.Value, 'f', -1, 64) + "</td>")
		body.WriteString(fmt.Sprintf("<td> %s </td>", DateOnFormat))
		body.WriteString(fmt.Sprintf("<td> %s </td>", ExpireOnFormat))
		body.WriteString("</tr>\n")
	}

	m.Set("BodyMessage", body.String())

	if err := m.SendMessage(); err != nil {
		return err
	}
	return nil
}

func AddICoupon(iCouponIn dto.ICouponIn, adminName string) (fiber.Map, int) {

	var iCoupon models.ICoupon
	var iCoupons []dto.SendCoupon

	result, email := repositories.GetUserEmailByDistribID(iCouponIn.DistribID)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return fiber.Map{"data": "No emailID Found"}, http.StatusNotFound
	}

	for _, Coupon := range iCouponIn.Coupons {
		for i := 0; i < int(Coupon.Quantity); i++ {
			hexVID := generateUniqueHexCode(10)
			hexPin := generateUniqueHexCode(10)
			iCoupon = models.ICoupon{
				VID:       hexVID,
				Pin:       hexPin,
				DistribID: iCouponIn.DistribID,
				DateOn:    time.Now(),
				TxDetail:  iCouponIn.TxDetail,
				ExpiresOn: time.Now().AddDate(0, 6, 0),
				Value:     Coupon.Value,
				AdminName: adminName,
			}

			SendCoupon := dto.SendCoupon{
				VID:       iCoupon.VID,
				Pin:       iCoupon.Pin,
				Value:     iCoupon.Value,
				DateOn:    iCoupon.DateOn,
				ExpiresOn: iCoupon.ExpiresOn,
			}
			iCoupons = append(iCoupons, SendCoupon)
			repositories.SaveICoupon(iCoupon)
		}
	}

	err := SendICouponMail(email, iCoupons)
	if err != nil {
		fmt.Print(err)
		return fiber.Map{"error1": err}, http.StatusInternalServerError
	}

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

	active, _ := repositories.ValidateICoupon(payload.VID, payload.Pin, distribID)
	if !active {
		return fiber.Map{"data": "Icoupon expired"}, http.StatusOK
	}
	balance, result := repositories.GetICouponBalance(payload.VID, payload.Pin, distribID)

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
