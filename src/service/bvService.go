package service

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBvHistory(payload dto.BvHistoryIn, tx *gorm.DB) (fiber.Map, int) {

	if payload.TransType != "" {

		BvHistory, err := repositories.GetBvHistoryByTransType(payload.DistribId, payload.TransType, payload.FromDate, payload.ToDate, tx)

		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetBvHistoryByTransType from GetBvHistory service fn ", err.Error())
			return fiber.Map{"data": "No BV Transaction History"}, fiber.StatusNotFound
		}
		return fiber.Map{"data": BvHistory}, fiber.StatusOK
	}

	BvHistory, err := repositories.GetAllBvHistory(payload.DistribId, tx)

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetAllBvHistory from GetBvHistory service fn ", err.Error())
		return fiber.Map{"data": "No Bv Transaction History"}, fiber.StatusNotFound
	}

	return fiber.Map{"data": BvHistory}, fiber.StatusOK
}

func GetBvDistributionTableByOrdeId(orderId string, distribId string, tx *gorm.DB) ([]dto.PlaceBv, error) {

	addedBv, err := repositories.GetAddedBvByOrderIdAndDistribId(orderId, distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetAddedBvByOrderIdAndDistribId from GetBvDistributionTableByOrdeId service fn ", err.Error())
		return addedBv, err
	}

	return addedBv, nil
}
