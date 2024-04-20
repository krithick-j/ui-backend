package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBvHistory(payload dto.BvHistoryIn) (fiber.Map, int) {

	if payload.TransType != "" {

		BvHistory, err := repositories.GetBvHistoryByTransType(payload.DistribId, payload.TransType, payload.FromDate, payload.ToDate)

		if err == gorm.ErrRecordNotFound {
			return fiber.Map{"data": "No BV Transaction History"}, fiber.StatusNotFound
		}
		return fiber.Map{"data": BvHistory}, fiber.StatusOK
	}

	BvHistory, err := repositories.GetAllBvHistory(payload.DistribId)

	if err == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No Bv Transaction History"}, fiber.StatusNotFound
	}

	return fiber.Map{"data": BvHistory}, fiber.StatusOK
}
