package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {

	router.Get("/:dist_id", controllers.GetUserByDistId)
}
