package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func RspRouter(router fiber.Router) {
	router.Get("/:distrib_id", controllers.GetTotalRspByDistribID)
	router.Get("/groupRsp/:distrib_id", controllers.GetGroupRspByDistribID)
}
