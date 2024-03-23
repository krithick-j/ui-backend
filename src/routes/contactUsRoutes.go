package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ContactCenter(router fiber.Router) {

	router.Post("/submitContactUsQuery", controllers.SubmitContactUsQuery)
	router.Post("/enquiryType", controllers.AddEnquiryType)
	router.Get("/allEnquiryType", controllers.GetAllEnquiryType)	
	
}
