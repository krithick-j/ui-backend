package controllers

import (
	"fmt"
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
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

func GetProductsByCategoryID(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")
	productType := c.Query("product_type")
	res, status := service.GetProductByCategoryID(categoryId, productType)

	return c.Status(status).JSON(res)
}

func AddToCartController(c *fiber.Ctx) error {
	var request dto.CartItemIn
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, status := service.AddToCart(request)

	return c.Status(status).JSON(res)
}

func GetProductsById(c *fiber.Ctx) error {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No IDs provided"})
	}

	res, _ := service.GetProductsByIds(request.IDs)
	return c.Status(http.StatusOK).JSON(res)

}

func GetEpProductsByCategoryId(c *fiber.Ctx) error {

	categoryId := c.Params("category_id")
	res, status := service.GetEpProductsByCategoryId(categoryId)

	return c.Status(status).JSON(res)
}

func GetCartProductsByUserId(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	res, status := service.GetCartProductsByDistribId(distrib_id)
	return c.Status(status).JSON(res)
}

func DeleteCartProduct(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	product_id := c.Query("product_id")
	res, status := service.DeleteCartProduct(distrib_id, product_id)
	return c.Status(status).JSON(res)
}

func DeleteAllCartProduct(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	res, status := service.DeleteAllCartProduct(distrib_id)
	return c.Status(status).JSON(res)
}

func EditCartProducts(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	product_id := c.Query("product_id")

	var payload models.CartItem
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, status := service.EditCartProducts(payload, distrib_id, product_id)
	return c.Status(status).JSON(res)
}

func CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err.Error())
	}

	res, status := service.CreateProduct(c, form)
	return c.Status(status).JSON(res)
}

func GetOrderDetails(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	res, status := service.GetOrderDetails(distrib_id)
	if status == http.StatusInternalServerError {
		return c.Status(status).JSON(fiber.Map{"error": res})
	}

	return c.Status(status).JSON(fiber.Map{"data": res})
}
