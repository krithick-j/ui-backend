package service

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn, orderDetails dto.OrderDetailsOut) (fiber.Map, int) {

	//Generating unique Order ID
	orderId := generateUniqueHexCode(10)

	//Placing Order
	//1.Adding in Header Table
	OrderHeaderObj := &models.OrdersHeader{
		DistribId: orderDetails.DistribId,
		Place:     OrderIn.Place,
		OrderId:   orderId,
	}
	repositories.SaveOrderHeader(OrderHeaderObj)

	//2. Adding in Footer Table
	OrderFooterObj := &models.OrderFooter{
		OrdersHeaderID: OrderHeaderObj.ID,
		SubTotal:       orderDetails.SubTotal,
		TotalSandH:     orderDetails.TotalSandH,
		TotalAmount:    OrderIn.TotalAmount, //total amount after applying coupon
		TotalQuantity:  orderDetails.TotalQuantity,
		TotalBV:        orderDetails.TotalBV,
		ContactName:    orderDetails.DeliveryAddress.ContactName,
		ContactEmail:   orderDetails.DeliveryAddress.ContactEmail,
		Address:        orderDetails.DeliveryAddress.Address,
		City:           orderDetails.DeliveryAddress.City,
		District:       orderDetails.DeliveryAddress.District,
		State:          orderDetails.DeliveryAddress.State,
		ZipCode:        orderDetails.DeliveryAddress.ZipCode,
		Country:        orderDetails.DeliveryAddress.Country,
		HomePhoneNo:    orderDetails.DeliveryAddress.HomePhoneNo,
		MobilePhoneNo:  orderDetails.DeliveryAddress.MobilePhoneNo,
	}
	repositories.SaveOrderFooter(OrderFooterObj)

	// 3. Adding in Liner Table
	for _, product := range orderDetails.Items {
		OrderLinerObj := &models.OrdersLiner{
			OrdersHeaderID: OrderHeaderObj.ID,
			Name:           product.Name,
			Quantity:       product.Quantity,
			UnitPrice:      product.UnitPrice,
			BV:             product.BV,
			SubTotal:       product.SubTotal,
			SandH:          product.SandH,
			OrderFooterID:  OrderFooterObj.ID,
		}
		repositories.SaveOrderLiner(OrderLinerObj)
	}

	// Close Coupon if coupon balance is 0
	for _, orderCoupon := range OrderIn.AppliedCoupons {
		totalValue := repositories.GetICouponValue(orderCoupon.VID, orderCoupon.Pin)
		balance, result := repositories.GetICouponBalance(orderCoupon.VID, orderCoupon.Pin, orderDetails.DistribId)
		if result.Error != nil {
			return fiber.Map{"error": result.Error}, http.StatusInternalServerError
		}

		//Recording in ICoupon Transaction table
		tx := models.ICouponTransaction{
			OrderId:        orderId,
			DistribId:      orderDetails.DistribId,
			VID:            orderCoupon.VID,
			Pin:            orderCoupon.Pin,
			TotalValue:     totalValue,
			AmountDetected: orderCoupon.AmountDetected,
			Balance:        balance - orderCoupon.AmountDetected,
		}

		//Closing coupon if balance is over
		if tx.Balance == 0 {
			repositories.CloseCoupon(orderCoupon.VID)
		}
		repositories.RecordICouponTx(tx)

	}

	// //Adding BV to the tree
	UpdateBvInTreeAfterPlaceOrder(orderId, orderDetails.DistribId, OrderIn.Place, orderDetails.TotalBV)

	return fiber.Map{"success": "Ordered Placed Successfully"}, http.StatusOK
}
