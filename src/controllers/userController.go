package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	data := struct {
		UserName string
		Password string
	}{}

	if err := c.BodyParser(&data); err != nil {
		configs.Log.Errorln("Error on Calling BodyParser data from Login Controller", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	res, status := service.LoginUser(data.UserName, data.Password, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on Committing tx from GetUserByDistId controller", err.Error())
	}
	return c.Status(status).JSON(res)
}

func GetUserByDistribId(c *fiber.Ctx) error {
	id := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetUserByDistribId(id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on Committing tx from GetUserByDistId controller", err.Error())
	}

	return c.Status(status).JSON(res)
}

func GetAllUsers(c *fiber.Ctx) error {

	tx := configs.DB.Begin()
	res, status := service.GetUsers(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on Committing tx from GetAllUsers controller", err.Error())
	}

	return c.Status(status).JSON(res)
}

func GetUserRank(c *fiber.Ctx) error {
	distribId := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetProfileDetails(distribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on Committing tx from GetAllUsers controller", err.Error())
	}

	return c.Status(status).JSON(res)
}

func GetUserTreeByDistId(c *fiber.Ctx) error {
	id := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetTreeUserByDistId(id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func UserRegistration(c *fiber.Ctx) error {
	user_in := dto.UserIn{}
	if err := c.BodyParser(&user_in); err != nil {
		configs.Log.Errorln("Error on Parsing user_in UserRegistration Controller", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	resp, status := service.RegisterUser(user_in, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(resp)

}

func GenerateConsentForm(c *fiber.Ctx) error {
	var payload struct {
		DistribId string `json:"distrib_id"`
	}
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on Parsing user_in UserRegistration Controller", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.GenerateConsentForm(payload.DistribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GenerateKYCDistribForm(c *fiber.Ctx) error {
	distribId := c.Params("distrib_id")
	tx := configs.DB.Begin()
	filename, status, err := service.GenerateDistributorForm(distribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	if err != nil {
		configs.Log.Errorln("Error on GenerateDistributorForm service from GenerateKYCDistribForm Controller", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return c.Status(status).JSON(fiber.Map{"filename": filename})
}

func EditUserByDistId(c *fiber.Ctx) error {

	DistribID := c.Params("distrib_id")

	var user_in models.User

	if err := c.BodyParser(&user_in); err != nil {
		configs.Log.Errorln("Error on parsing user_in from EditUserByDistId controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.EditUserByDistId(DistribID, user_in, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func NewReferrals(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetNewReferrals(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetTrackingCenters(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetTrackingCentersByDistribId(distrib_id, true, "", "", tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func UpdateUserPass(c *fiber.Ctx) error {

	var user_in dto.UserPassIn

	if err := c.BodyParser(&user_in); err != nil {
		configs.Log.Errorln("Error on parsing user_in from UpdateUserPass controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	res, status := service.UpdateUserPass(user_in, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetIDCard(c *fiber.Ctx) error {
	data := struct {
		DistribId string `json:"distrib_id"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		configs.Log.Errorln("Erron on parsing data from GetIDCard controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	tx := configs.DB.Begin()
	res, status := service.GenerateIDCard(data.DistribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetMediaFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	c.Status(fiber.StatusOK).SendFile("assets/" + filename)
	return nil
}

func SendEmailCode(c *fiber.Ctx) error {

	data := struct {
		Email string `json:"email"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		configs.Log.Errorln("Error on Parsing data from SendEmailCode controller fn")
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	tx := configs.DB.Begin()
	res, status := service.SendEmailCode(data.Email, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	c.Status(status).JSON(res)
	return nil
}

func VerifyEmailCode(c *fiber.Ctx) error {

	payload := struct {
		OTP   string `json:"otp"`
		Email string `json:"email"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on calling BodyParser fn from VerifyEmailCode controllers fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	res, status := service.VerifyEmailCode(payload.OTP, payload.Email, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func SendPhoneCode(c *fiber.Ctx) error {

	data := struct {
		PhoneNo string `json:"phone_no"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		configs.Log.Errorln("Error on Parsing data from SendPhoneCode fn", err.Error())
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	tx := configs.DB.Begin()
	res, status := service.SendPhoneCode(c, data.PhoneNo, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func VerifyPhoneCode(c *fiber.Ctx) error {

	data := struct {
		OTP     string `json:"otp"`
		PhoneNo string `json:"phone_no"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		configs.Log.Errorln("Error on Parsing data from VerifyPhoneCode controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	tx := configs.DB.Begin()
	res, status := service.VerifyPhoneCode(data.OTP, data.PhoneNo, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func KycUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		configs.Log.Errorln("Error on calling MultipartForm fn from KycUpload controller fn", err.Error())
	}
	tx := configs.DB.Begin()
	res, status := service.KycUpload(c, form, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func ApproveKYC(c *fiber.Ctx) error {
	data := struct {
		DistribId string `json:"distrib_id"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		configs.Log.Errorln("Error on parsing data from ApproveKYC controller function", err.Error())
		c.Status(fiber.StatusBadRequest).SendString("{\"error\":\"Bad Request\"}")
	}
	tx := configs.DB.Begin()
	res, status := service.ApproveKYC(data.DistribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}
