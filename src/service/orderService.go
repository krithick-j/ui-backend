package service

import (
	"net/http"
	"ui-back-end/configs"
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
	total, status := GetOrderDetails(OrderIn.DistribId)
	tx := configs.DB.Begin()

	if status != http.StatusOK {
		return fiber.Map{"data": "Something gone wrong"}, http.StatusInternalServerError
	}

	productType := total.Items[0].ProductType

	//Placing Order
	//1.Adding in Header Table
	OrderHeaderObj := &models.OrdersHeader{
		DistribId:      OrderIn.DistribId,
		OrderId:        orderId,
		SubTotal:       total.SubTotal,
		TotalSandH:     total.TotalSandH,
		TotalAmount:    OrderIn.TotalAmount, //total amount after applying coupon
		TotalQuantity:  total.TotalQuantity,
		TotalTypeValue: total.TotalTypeValue,
		ProductType:    productType,
		ContactName:    total.DeliveryAddress.ContactName,
		ContactEmail:   total.DeliveryAddress.ContactEmail,
		Address:        total.DeliveryAddress.Address,
		City:           total.DeliveryAddress.City,
		District:       total.DeliveryAddress.District,
		State:          total.DeliveryAddress.State,
		ZipCode:        total.DeliveryAddress.ZipCode,
		Country:        total.DeliveryAddress.Country,
		HomePhoneNo:    total.DeliveryAddress.HomePhoneNo,
		MobilePhoneNo:  total.DeliveryAddress.MobilePhoneNo,
	}
	err := repositories.SaveOrderHeader(OrderHeaderObj)

	if err != nil {
		tx.Rollback()
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	// 3. Adding in Liner Table
	for _, product := range total.Items {
		OrderLinerObj := &models.OrdersLiner{
			OrdersHeaderID: OrderHeaderObj.ID,
			Name:           product.Name,
			Quantity:       product.Quantity,
			UnitPrice:      product.UnitPrice,
			ProductType:    product.ProductType,
			TypeValue:      product.TypeValue,
			SubTotal:       product.SubTotal,
			SandH:          product.SandH,
		}
		err = repositories.SaveOrderLiner(OrderLinerObj)
		if err != nil {
			tx.Rollback()
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
	}

	//validate coupon balance
	// Close Coupon if coupon balance is 0
	if OrderIn.AppliedCoupons != nil && productType != "ep" {
		for _, orderCoupon := range OrderIn.AppliedCoupons {
			totalValue := repositories.GetICouponValue(orderCoupon.VID, orderCoupon.Pin)
			balance, result := repositories.GetICouponBalance(orderCoupon.VID)
			if result.Error != nil {
				return fiber.Map{"error": result.Error.Error()}, http.StatusInternalServerError
			}

			//Recording in ICoupon Transaction table
			ICouponObj := models.ICouponTransaction{
				OrderId:        orderId,
				DistribId:      OrderIn.DistribId,
				VID:            orderCoupon.VID,
				Pin:            orderCoupon.Pin,
				TotalValue:     totalValue,
				AmountDetected: orderCoupon.AmountDetected,
				Balance:        balance - orderCoupon.AmountDetected,
			}

			//Updating balance in icoupons Transaction table
			_, res := repositories.UpdateBalanceInICoupons(orderCoupon.VID, ICouponObj.Balance)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}

			//Closing coupon if balance is over
			if ICouponObj.Balance == 0 {
				res = repositories.CloseCoupon(orderCoupon.VID)
				if res.Error != nil {
					tx.Rollback()
					return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
				}
			}
			res = repositories.SaveICouponTx(ICouponObj)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		}
		//if len(bv)=0 then rsp transaction if not bv transaction
		if len(OrderIn.PlaceBvs) != 0 && productType == "bv" {
			//insert directbv in rsptransaction
			res := repositories.AddDirectBvTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
			//bv transaction insert is done in UpdateCurrentPlaceValues function
			UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId, tx)
			//UpdateTreePlaceValuesByDistribId(OrderIn.DistribId, OrderIn.PlaceBvs)
		} else if productType == "rsp" {
			//save rsp transaction
			repositories.AddRspTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
		}
	} else if productType == "ep" {
		totalEp, res := repositories.GetEpBalance(total.DistribId)
		if res.Error != nil {
			tx.Rollback()
			return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
		}
		//balance should be atleast gte to the order of buying product
		if totalEp >= total.TotalTypeValue {
			res = repositories.SaveEpTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		} else {
			return fiber.Map{"data": "Insufficient Balance!"}, fiber.StatusForbidden
		}
	}

	//Cleaning cart after buying
	repositories.DeleteAllCartProduct(OrderIn.DistribId)

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
