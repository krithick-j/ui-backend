package controllers

import (
	"net/http"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetProductsController(c *fiber.Ctx) error {
	res := service.GetProducts()
	return c.Status(http.StatusOK).JSON(res)
}

func GetProductCategoriesController(c *fiber.Ctx) error {
	res := service.GetProductCategories()
	return c.Status(http.StatusOK).JSON(res)
}
