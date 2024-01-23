package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func GrRouter(router fiber.Router) {

	router.Get("/allGrVisual", controllers.GetAllGrVisual)
	router.Get("/allGrVisualByDate", controllers.GetAllGrVisualByDate)

}
