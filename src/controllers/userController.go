package controllers

import (
	"fmt"
	"net/http"
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
	c.BodyParser(&data)
	res, status := service.LoginUser(data.UserName, data.Password)
	return c.Status(status).JSON(res)
}

func GetUserByDistId(c *fiber.Ctx) error {
	id := c.Params("dist_id")

	res, status := service.GetUserByDistId(id)

	return c.Status(status).JSON(res)
}

func GetAllUsers(c *fiber.Ctx) error {
	res := service.GetUsers()
	return c.Status(http.StatusOK).JSON(res)
}

func GetUserTreeByDistId(c *fiber.Ctx) error {
	id := c.Params("distrib_id")
	res := service.GetTreeUserByDistId(id)
	return c.Status(http.StatusOK).JSON(res)
}

func UserRegistration(c *fiber.Ctx) error {
	user_in := dto.UserIn{}
	c.BodyParser(&user_in)
	resp, err := service.RegisterUser(user_in)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(resp)
	} else {
		return c.Status(http.StatusCreated).JSON(resp)
	}
}

func EditUserByDistId(c *fiber.Ctx) error {

	DistribID := c.Params("distrib_id")

	var user_in models.User

	if err := c.BodyParser(&user_in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	res, status := service.EditUserByDistId(DistribID, user_in)

	return c.Status(status).JSON(res)
}

func NewReferrals(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	res, status := service.GetNewReferrals(distrib_id)
	return c.Status(status).JSON(res)
}

func GetTrackingCenters(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	res, status := service.GetTrackingCentersByDistribId(distrib_id)
	return c.Status(status).JSON(res)
}

func UpdateUserPass(c *fiber.Ctx) error {

	var user_in dto.UserPassIn

	if err := c.BodyParser(&user_in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	res, status := service.UpdateUserPass(user_in)

	return c.Status(status).JSON(res)
}

func GetIDCard(c *fiber.Ctx) error {
	data := struct {
		DistribId string `json:"distrib_id"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	fmt.Println(data.DistribId)
	service.GenerateIDCard(data.DistribId)
	return c.Status(fiber.StatusCreated).JSON("{msg:success}")
}

func GetMediaFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	print(filename)
	c.Status(fiber.StatusOK).SendFile("assets/" + filename)
	return nil
}

func SendEmailCode(c *fiber.Ctx) error {

	data := struct {
		Email string `json:"email"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	service.SendEmailCode(data.Email)
	c.Status(fiber.StatusCreated).SendString("Created")
	return nil
}

func VerifyEmailCode(c *fiber.Ctx) error {

	data := struct {
		OTP   string `json:"otp"`
		Email string `json:"email"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	err = service.VerifyEmailCode(data.OTP, data.Email)
	if err == nil {
		c.Status(fiber.StatusOK).SendString("Ok")
	} else {
		c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	return nil
}

func SendPhoneCode(c *fiber.Ctx) error {

	data := struct {
		PhoneNo string `json:"phone_no"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	service.SendPhoneCode(data.PhoneNo)
	c.Status(fiber.StatusCreated).SendString("Created")
	return nil
}

func VerifyPhoneCode(c *fiber.Ctx) error {

	data := struct {
		OTP     string `json:"otp"`
		PhoneNo string `json:"phone_no"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	fmt.Println("data", data)
	err = service.VerifyPhoneCode(data.OTP, data.PhoneNo)
	if err == nil {
		c.Status(fiber.StatusOK).SendString("Ok")
	} else {
		c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}
	return nil
}

func KycUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err.Error())
	}
	service.KycUpload(c, form)
	fmt.Println("Processing Files")
	c.Status(fiber.StatusAccepted).SendString("Accepted")
	return nil
}

func ApproveKYC(c *fiber.Ctx) error {
	data := struct {
		DistribId string `json:"distrib_id"`
	}{}
	err := c.BodyParser(&data)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString("{\"error\":\"Bad Request\"}")
	}
	err = service.ApproveKYC(data.DistribId)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString("{\"error\":\"Bad Request\"}")
	}
	return nil
}
