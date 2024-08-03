package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func RspRouter(router fiber.Router) {
	router.Get("personalRsp/:distrib_id", controllers.GetPersonalRspByDistribID)
	router.Get("/groupRsp/:distrib_id", controllers.GetGroupRspByDistribID)
	router.Get("/directBv/:distrib_id", controllers.GetDirectBvByDistribID)
}
