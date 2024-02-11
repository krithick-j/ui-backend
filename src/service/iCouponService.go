package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
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

func SendICouponMail(toMail string, VID string, Pin string) error {
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

	m.Set("BodyMessage", fmt.Sprintf("This is a noreply email. Your ICoupon VID is %s and Pin is %s", VID, Pin))

	if err := m.SendMessage(); err != nil {
		return err
	}
	return nil
}

func AddICoupon(iCouponIn dto.ICouponIn, adminName string) (fiber.Map, int) {

	var iCoupon models.ICoupon
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
				DateOn:    Coupon.DateOn,
				TxDetail:  iCouponIn.TxDetail,
				ExpiresOn: Coupon.ExpiresOn,
				Value:     Coupon.Value,
				AdminName: adminName,
			}
			err := SendICouponMail(email, hexVID, hexPin)
			if err != nil {
				fmt.Print(err)
				return fiber.Map{"error1": err}, http.StatusInternalServerError
			}
			repositories.SaveICoupon(iCoupon)
		}
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

	var iCoupon models.ICoupon

	iCoupon, result := repositories.GetICoupon(payload.VID, payload.Pin, iCoupon, distribID)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon exists"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": iCoupon}, http.StatusOK
}
