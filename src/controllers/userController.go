package controllers

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

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
	id := c.Params("dist_id")
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
