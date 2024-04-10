package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBvHistory(payload dto.BvHistoryIn) (fiber.Map, int) {

	if payload.HistoryType != "" {
		iCouponHistory, err := repositories.GetBvHistory(payload.DistribId, payload.HistoryType)

		if err == gorm.ErrRecordNotFound {
			return fiber.Map{"data": "No ICoupon Transaction History"}, fiber.StatusNotFound
		}
		return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
	}

	iCouponHistory, err := repositories.GetAllBvHistory(payload.DistribId)

	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No ICoupon Transaction History"}, fiber.StatusNotFound
	}

	return fiber.Map{"data": iCouponHistory}, fiber.StatusOK
}
