package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func CpaRouter(router fiber.Router) {

	router.Get("/settlementCatalog", controllers.GetSettlementCatalog)

}
