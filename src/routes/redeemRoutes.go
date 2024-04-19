package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func RedeemRouter(router fiber.Router) {
	router.Get("/:distrib_id", controllers.GetEpBalanceByDistribId)
}
