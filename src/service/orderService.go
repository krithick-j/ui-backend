package service

import (
	"fmt"
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn) (fiber.Map, int) {

	//Generating unique Order ID
	orderId := GenerateUniqueHexCode(10)

	//sum product value
	res, status := GetOrderDetails(OrderIn.DistribId)
	if status != http.StatusOK {
		return fiber.Map{"data": "Something gone wrong"}, http.StatusInternalServerError
	}

	//Placing Order
	//1.Adding in Header Table
	OrderHeaderObj := &models.OrdersHeader{
		DistribId:     OrderIn.DistribId,
		OrderId:       orderId,
		SubTotal:      res.SubTotal,
		TotalSandH:    res.TotalSandH,
		TotalAmount:   OrderIn.TotalAmount, //total amount after applying coupon
		TotalQuantity: res.TotalQuantity,
		TotalBV:       res.TotalBV,
		ContactName:   res.DeliveryAddress.ContactName,
		ContactEmail:  res.DeliveryAddress.ContactEmail,
		Address:       res.DeliveryAddress.Address,
		City:          res.DeliveryAddress.City,
		District:      res.DeliveryAddress.District,
		State:         res.DeliveryAddress.State,
		ZipCode:       res.DeliveryAddress.ZipCode,
		Country:       res.DeliveryAddress.Country,
		HomePhoneNo:   res.DeliveryAddress.HomePhoneNo,
		MobilePhoneNo: res.DeliveryAddress.MobilePhoneNo,
	}
	repositories.SaveOrderHeader(OrderHeaderObj)

	// 3. Adding in Liner Table
	for _, product := range res.Items {
		OrderLinerObj := &models.OrdersLiner{
			OrdersHeaderID: OrderHeaderObj.ID,
			Name:           product.Name,
			Quantity:       product.Quantity,
			UnitPrice:      product.UnitPrice,
			BV:             product.BV,
			RSP:            product.Rsp,
			EP:             product.Ep,
			SubTotal:       product.SubTotal,
			SandH:          product.SandH,
		}
		repositories.SaveOrderLiner(OrderLinerObj)
	}
	
	//validate coupon balance
	// Close Coupon if coupon balance is 0
	if OrderIn.AppliedCoupons != nil && OrderIn.OrderType != "ep" {
		for _, orderCoupon := range OrderIn.AppliedCoupons {
			totalValue := repositories.GetICouponValue(orderCoupon.VID, orderCoupon.Pin)
			balance, result := repositories.GetICouponBalance(orderCoupon.VID)
			if result.Error != nil {
				return fiber.Map{"error": result.Error}, http.StatusInternalServerError
			}

			//Recording in ICoupon Transaction table
			tx := models.ICouponTransaction{
				OrderId:        orderId,
				DistribId:      OrderIn.DistribId,
				VID:            orderCoupon.VID,
				Pin:            orderCoupon.Pin,
				TotalValue:     totalValue,
				AmountDetected: orderCoupon.AmountDetected,
				Balance:        balance - orderCoupon.AmountDetected,
			}

			//Updating balance in icoupons Transaction table
			repositories.UpdateBalanceInICoupons(orderCoupon.VID, tx.Balance)

			//Closing coupon if balance is over
			if tx.Balance == 0 {
				repositories.CloseCoupon(orderCoupon.VID)
			}
			repositories.RecordICouponTx(tx)

		}
		//if len(bv)=0 then rsp transaction if not bv transaction
		if len(OrderIn.PlaceBvs) != 0 {
			//insert directbv in rsptransaction
			repositories.AddDirectBvTx(OrderIn.DistribId, orderId, res.TotalBV)
			UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId)
			//UpdateTreePlaceValuesByDistribId(OrderIn.DistribId, OrderIn.PlaceBvs)
		} else {
			//save rsp transaction
			fmt.Println("total rsp from service", res.TotalRsp)
			repositories.AddRspTx(OrderIn.DistribId, orderId, res.TotalRsp)
		}
	}

	//Cleaning cart after buying
	var cartItems []models.CartItem
	repositories.DeleteAllCartProduct(OrderIn.DistribId, cartItems)

	return fiber.Map{"success": "Ordered Placed Successfully"}, http.StatusOK

}

func GetOrdersByDistribId(distrib_id string) (fiber.Map, int) {

	var order []models.OrdersHeader
	var result *gorm.DB

	order, result = repositories.GetOrderByDistribId(distrib_id, order)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"data": order}, http.StatusOK
}
