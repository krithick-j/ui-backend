package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {

	//User Crud
	router.Get("/", controllers.GetAllUsers)
	router.Get("/:dist_id", controllers.GetUserByDistId)
	router.Put("/:distrib_id", controllers.EditUserByDistId)
	//User Chain
	router.Get("/tree/:distrib_id", controllers.GetUserTreeByDistId)
	router.Get("/newReferral/:distrib_id", controllers.NewReferrals)
	router.Get("/allReferrals", controllers.AllReferrals)
	router.Get("/getTrackingCenters/:distrib_id", controllers.GetTrackingCenters)
}
