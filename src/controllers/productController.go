package controllers

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetProductsController(c *fiber.Ctx) error {
	tx := configs.DB.Begin()
	res, status := service.GetProducts(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetProductCategoriesController(c *fiber.Ctx) error {
	tx := configs.DB.Begin()
	res, status := service.GetProductCategories(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetProductsByCategoryID(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")
	productType := c.Query("product_type")
	tx := configs.DB.Begin()
	res, status := service.GetProductByCategoryID(categoryId, productType, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func AddToCartController(c *fiber.Ctx) error {
	var request dto.CartItemIn
	if err := c.BodyParser(&request); err != nil {
		configs.Log.Errorln("Error on parsing request from AddToCartController controller fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	tx := configs.DB.Begin()
	res, status := service.AddToCart(request, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetProductsById(c *fiber.Ctx) error {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.BodyParser(&request); err != nil {
		configs.Log.Errorln("Error on parsing request from GetProductsById controller fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No IDs provided"})
	}
	tx := configs.DB.Begin()
	res, status := service.GetProductsByIds(request.IDs, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetEpProductsByCategoryId(c *fiber.Ctx) error {

	categoryId := c.Params("category_id")
	tx := configs.DB.Begin()
	res, status := service.GetEpProductsByCategoryId(categoryId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetCartProductsByUserId(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetCartProductsByDistribId(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func DeleteCartProduct(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	product_id := c.Query("product_id")
	tx := configs.DB.Begin()
	res, status := service.DeleteCartProduct(distrib_id, product_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func DeleteAllCartProduct(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.DeleteAllCartProduct(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func EditCartProducts(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	product_id := c.Query("product_id")

	var payload models.CartItem
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from EditCartProducts controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	tx := configs.DB.Begin()
	res, status := service.EditCartProducts(payload, distrib_id, product_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		configs.Log.Errorln("Error on parsing multipartForm from CreateProduct", err.Error())
		fmt.Println(err.Error())
	}
	tx := configs.DB.Begin()
	res, status := service.CreateProduct(c, form, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

// func EditProduct(c *fiber.Ctx) error {
// 	distrib_id := c.Query("distrib_id")
// 	product_id := c.Query("product_id")

// 	var payload models.Product
// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
// 	}

// 	res, status := service.EditProduct(payload, distrib_id, product_id)
// 	return c.Status(status).JSON(res)
// }

func GetOrderDetails(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetOrderDetails(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(fiber.Map{"data": res})
}
