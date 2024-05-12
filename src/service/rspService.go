package service

import (
	"fmt"
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTotalRspByDistribId(DistribID string) (fiber.Map, int) {
	var totalRsp int = 0

	rsp, result := repositories.GetAllRspByDistribId(DistribID)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "No transaction Record Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	for _, rspValue := range rsp {
		totalRsp += rspValue
	}

	Output := dto.TotalRspOut{
		TotalRsp: totalRsp,
	}

	return fiber.Map{"data": Output}, http.StatusOK

}
