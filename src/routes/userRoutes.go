package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {

	//User Crud
	router.Post("/idCard", controllers.GetIDCard)
	router.Post("/distribApplicationForm", controllers.GenerateConsentForm)
	router.Get("/", controllers.GetAllUsers)
	router.Get("/:dist_id", controllers.GetUserByDistId)
	router.Put("/updateUser/:distrib_id", controllers.EditUserByDistId)
	router.Put("/changePass", controllers.UpdateUserPass)
	router.Get("/acknowledgementLetter/:distribId", controllers.GetMediaFile)
	//User Chain
	router.Post("/kyc-upload", controllers.KycUpload)
	router.Get("/tree/:distrib_id", controllers.GetUserTreeByDistId)
	router.Get("/newReferral/:distrib_id", controllers.NewReferrals)
	router.Get("/allReferrals", controllers.AllReferrals)
	router.Get("/getTrackingCenters/:distrib_id", controllers.GetTrackingCenters)
	router.Post("/kycapprove", controllers.ApproveKYC)
}
