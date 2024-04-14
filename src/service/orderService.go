package service

import (
	"fmt"
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
	totalOrderAmount := total.TotalAmount

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
		TotalAmount:    total.TotalAmount, //total amount before applying coupon
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

	var totalICouponBalance float64
	if OrderIn.AppliedCoupons != nil && productType != "ep" {
		for _, orderCoupon := range OrderIn.AppliedCoupons {
			balance, result := repositories.GetICouponBalance(orderCoupon.VID)
			if result.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": result.Error.Error()}, http.StatusInternalServerError
			}

			if balance == 0 {
				repositories.CloseCoupon(orderCoupon.VID)
				tx.Rollback()
				return fiber.Map{"error": "ICoupon used already!"}, fiber.StatusPaymentRequired
			}
			totalICouponBalance += balance
		}

		for _, orderCoupon := range OrderIn.AppliedCoupons {
			balance, result := repositories.GetICouponBalance(orderCoupon.VID)
			if result.Error != nil {
				return fiber.Map{"error": result.Error.Error()}, http.StatusInternalServerError
			}
			fmt.Println(totalOrderAmount)
			if totalOrderAmount >= balance {
				totalOrderAmount = totalOrderAmount - balance

				// Updating in ICoupon Transaction table
				ICouponObj := models.ICouponTransaction{
					DistribId: OrderIn.DistribId,
					VID:       orderCoupon.VID,
					Reference: orderId,
					Value:     -balance, //using the full balance of ICoupon
				}
				fmt.Println("Balance ", balance)
				err := repositories.SaveICouponTx(ICouponObj)
				if err.Error != nil {
					tx.Rollback()
					return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError
				}
				//Balance will be zero after using full coupon, so active set to false
				res := repositories.CloseCoupon(orderCoupon.VID)
				if res.Error != nil {
					tx.Rollback()
					return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
				}

			} else {
				ICouponObj := models.ICouponTransaction{
					DistribId: OrderIn.DistribId,
					VID:       orderCoupon.VID,
					Value:     -totalOrderAmount, //The order amount is less than Icoupon Balance, so icoupon have remaining balance
				}
				fmt.Println("Balance from else block")
				err := repositories.SaveICouponTx(ICouponObj)
				if err.Error != nil {
					tx.Rollback()
					return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError
				}
				break //No need to loop again, since the order amount is satisfied with the coupon
			}
		}

		if totalICouponBalance < total.TotalAmount {
			tx.Rollback()
			return fiber.Map{"error": "Insufficient Balance! Transaction Failed"}, fiber.StatusInternalServerError
		}

		//if len(bv)=0 then rsp transaction if not bv transaction
		if len(OrderIn.PlaceBvs) != 0 && productType == "bv" {
			//insert directbv in rsptransaction
			res := repositories.AddDirectBvTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
			fmt.Println("-----------------------------Tree starts from here-----------------------------------------")
			//bv transaction insert is done in UpdateCurrentPlaceValues function
			UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId, tx, total.TotalTypeValue)
			fmt.Println("----------------------------Tree ends here--------------------------------------------------")
			//UpdateTreePlaceValuesByDistribId(OrderIn.DistribId, OrderIn.PlaceBvs)
		} else if productType == "rsp" {
			//save rsp transaction
			res := repositories.AddRspTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
			if res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		} else {
			tx.Rollback()
			return fiber.Map{"error": "Product type does not match"}, fiber.StatusInternalServerError
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
	_, res := repositories.DeleteAllCartProduct(OrderIn.DistribId)
	if res.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}

	res = tx.Commit()
	if res.Error != nil {
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}
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
