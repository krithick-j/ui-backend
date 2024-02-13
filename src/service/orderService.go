package service

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn, distribId string) (fiber.Map, int) {
	var coupon models.ICoupon

	//Placing Order
	for _, product := range OrderIn.Products {
		orderProduct := models.Orders{
			DistribId: distribId,
			ProductID: product.Productid,
		}
		repositories.SaveToOrders(orderProduct)
	}

	// Deleting Coupons after ordering product
	for _, Ordercoupon := range OrderIn.Coupons {
		coupon = models.ICoupon{
			VID: Ordercoupon.VID,
		}
		repositories.DeleteICoupon(coupon)
	}

	return fiber.Map{"success": "Added to Cart Successfully"}, http.StatusOK
}
