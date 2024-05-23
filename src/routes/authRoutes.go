package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router) {
	router.Post("/auth/login", controllers.Login)
	router.Post("/auth/register", controllers.UserRegistration)
	router.Post("/auth/emailcode", controllers.SendEmailCode)
	router.Put("/auth/emailcode", controllers.VerifyEmailCode)
	router.Post("/auth/phonecode", controllers.SendPhoneCode)
	router.Put("/auth/phonecode", controllers.VerifyPhoneCode)
}
