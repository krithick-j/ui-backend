package service

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func GetProducts() fiber.Map {
	var products []models.Product
	var result *gorm.DB
	products, result = repositories.GetAllProducts(products)
	println(result)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": products}
}

func GetProductCategories() fiber.Map {
	var productCategories []models.ProductCategory
	var result *gorm.DB
	productCategories, result = repositories.GetAllProductCategories(productCategories)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": productCategories}
}

func GetProductByCategoryID(category_id string) fiber.Map {

	var product []models.Product
	var result *gorm.DB

	product, result = repositories.GetAllProductByCategoryID(category_id, product)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": product}
}

func GetProductsByIds(ids []uint) ([]models.Product, error) {
	var products []models.Product
	products, result := repositories.GetAllProductByIDs(ids, products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func AddToCart(request dto.CartItemIn) fiber.Map {

	for _, item := range request.Items {
		product := models.CartItem{
			UserID:    request.UserID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		repositories.SaveToCart(product)
	}

	return fiber.Map{"success": "Added to Cart Successfully"}
}

func GetCartProductsByUserId(user_id string) fiber.Map {

	var products []dto.ProductsOut

	products, result := repositories.GetAllCartProductsByUserID(user_id, products)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}

	if result.RowsAffected == 0 {
		return fiber.Map{"data": "No Products in Cart"}
	}

	return fiber.Map{"data": products}
}

func DeleteCartProduct(user_id string, product_id string) (fiber.Map, int) {
	var cartItem models.CartItem
	cartItem, result := repositories.DeleteCartProduct(user_id, product_id, cartItem)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"success": "Product Not Found", "DeletedProduct": cartItem}, http.StatusNoContent
	}

	return fiber.Map{"success": "Product Deleted Successfully", "DeletedProduct": cartItem}, http.StatusOK
}

func EditCartProducts(payload models.CartItem, user_id string, product_id string) (fiber.Map, int) {

	cartItem, result := repositories.EditCartProducts(user_id, product_id, payload)

	if result.Error != nil {
		log.Info("Error saving user to the database:", result.Error)
		return fiber.Map{"error": result.Error}, http.StatusBadGateway
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"success": "Product Updated Successfully", "UpdatedProduct": cartItem}, http.StatusOK
}
