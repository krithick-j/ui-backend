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
	router.Put("/updateDp", controllers.UpdateDp)
	router.Get("/profile/:distrib_id", controllers.GetProfileDetails)
	router.Get("/:distrib_id", controllers.GetUserByDistribId)
	router.Put("/updateUser/:distrib_id", controllers.EditUserByDistId)
	router.Put("/changePass", controllers.UpdateUserPass)
	router.Put("/resetPass", controllers.ResetPass)
	router.Get("/acknowledgementLetter/:distrib_id", controllers.GetMediaFile)
	//User Chain
	router.Get("/tree/:distrib_id", controllers.GetUserTreeByDistId)
	router.Get("/newReferral/:distrib_id", controllers.NewReferrals)
	router.Get("/allReferrals", controllers.AllReferrals)
	router.Get("/getTrackingCenters/:distrib_id", controllers.GetTrackingCenters)

	//KYC routes
	router.Get("/kycdistribform/:distrib_id", controllers.GenerateKYCDistribForm)
	router.Post("/kyc-upload", controllers.KycUpload)
	router.Post("/kycapprove", controllers.ApproveKYC)
}
