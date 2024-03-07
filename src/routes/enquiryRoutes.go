package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func EnquiryRouter(router fiber.Router) {
	router.Post("/", controllers.AddEnquiryType)
}
