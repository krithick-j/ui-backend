package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func TestRouter(router fiber.Router) {
	router.Get("/:orderId", controllers.TestController)
}
