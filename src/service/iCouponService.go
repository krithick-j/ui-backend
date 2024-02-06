package service

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
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

func AddICoupon(iCouponIn dto.ICouponIn, adminName string) (fiber.Map, int) {

	var iCoupon models.ICoupon

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
			repositories.SaveICoupon(iCoupon)
		}
	}

	// if err := repositories.ProductSave(&newProduct); err != nil {
	// 	log.Info("Error saving user to the database:", err)
	// 	return fiber.Map{}
	// }

	// log.Info("Product uploaded successfully.")

	return fiber.Map{"data": "ICoupons added successfully"}, http.StatusCreated
}
