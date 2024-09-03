package service

import (
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SocialConnService(tx *gorm.DB) (fiber.Map, int) {
	data, err := repositories.GetAllSocialConn(tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllSocialConn", "SocialConnService", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(data, fiber.StatusOK)
}
