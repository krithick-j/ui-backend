package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn, distribId string) fiber.Map {

	

	for _, item := range request.Items {
		product := models.CartItem{
			DistribID: request.DistribID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		repositories.SaveToCart(product)
	}

	return fiber.Map{"success": "Added to Cart Successfully"}
}
