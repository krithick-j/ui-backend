package controllers

import (
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