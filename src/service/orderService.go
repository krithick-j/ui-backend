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
	if status != http.StatusOK {
		return fiber.Map{"data": "Something gone wrong"}, fiber.StatusInternalServerError
	}

	tx := configs.DB.Begin()

	totalOrderAmount := total.TotalAmount

	productType := total.Items[0].ProductType

	if res, status := handleProductHeaderAndLines(OrderIn, orderId, total, productType, tx); status != fiber.StatusOK {
		return fiber.Map{"error": res["error"]}, status
	}

	if res, status := handleProductType(OrderIn, productType, orderId, totalOrderAmount, tx, total); status != fiber.StatusOK {
		return fiber.Map{"error": res["error"]}, status
	}

	if commitRes := tx.Commit(); commitRes.Error != nil {
		return fiber.Map{"error": commitRes.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"success": "Ordered Placed Successfully"}, http.StatusOK
}

func handleProductType(OrderIn dto.PlaceOrderIn, productType string, orderId string, totalOrderAmount float64, tx *gorm.DB, total dto.OrderDetailsOut) (fiber.Map, int) {
	if OrderIn.AppliedCoupons != nil && productType != "ep" {

		if res, status := handlePlaceOrderICoupons(OrderIn.AppliedCoupons, OrderIn.DistribId, orderId, totalOrderAmount, tx); status != fiber.StatusOK {
			return fiber.Map{"error": res["error"]}, status
		}

		if productType == "bv" {

			if res, status := handleBvProduct(OrderIn, orderId, total, tx); status != fiber.StatusOK {
				return fiber.Map{"error": res["error"]}, fiber.StatusInternalServerError
			}

		} else if productType == "rsp" {

			if res := repositories.AddRspTx(OrderIn.DistribId, orderId, total.TotalTypeValue); res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		} else {
			tx.Rollback()
			return fiber.Map{"error": "Product type does not match"}, fiber.StatusInternalServerError
		}

	} else if productType == "ep" {

		if res, status := handleEpProduct(total, tx, OrderIn, orderId); status != fiber.StatusOK {
			return fiber.Map{"error": res["error"]}, fiber.StatusInternalServerError
		}
	}

	_, cartRes := repositories.DeleteAllCartProduct(OrderIn.DistribId)
	if cartRes.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": cartRes.Error.Error()}, fiber.StatusInternalServerError
	}
	return nil, fiber.StatusOK
}

func handleEpProduct(total dto.OrderDetailsOut, tx *gorm.DB, OrderIn dto.PlaceOrderIn, orderId string) (fiber.Map, int) {
	totalEp, res := repositories.GetEpBalance(total.DistribId)
	if res.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}

	if totalEp >= total.TotalTypeValue {
		res = repositories.SaveEpTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
		if res.Error != nil {
			tx.Rollback()
			return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
		}
	} else {
		return fiber.Map{"data": "Insufficient Balance!"}, fiber.StatusForbidden
	}
	return nil, fiber.StatusOK
}

func handleBvProduct(OrderIn dto.PlaceOrderIn, orderId string, total dto.OrderDetailsOut, tx *gorm.DB) (fiber.Map, int) {

	res := repositories.AddDirectBvTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
	if res.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}

	msg, status := UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId, tx, total.TotalTypeValue)
	if status == fiber.StatusInternalServerError {
		return fiber.Map{"error": msg}, status
	}

	Dcmessage, status := SaveDirectCommissionTransaction(OrderIn.DistribId, total.TotalTypeValue, orderId)
	if status != 200 {
		return fiber.Map{"error": Dcmessage}, status
	}

	return nil, fiber.StatusOK
}

func handleProductHeaderAndLines(OrderIn dto.PlaceOrderIn, orderId string, total dto.OrderDetailsOut, productType string, tx *gorm.DB) (fiber.Map, int) {

	OrderHeaderObj := &models.OrdersHeader{
		DistribId:      OrderIn.DistribId,
		OrderId:        orderId,
		SubTotal:       total.SubTotal,
		TotalSandH:     total.TotalSandH,
		TotalAmount:    total.TotalAmount,
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
	return fiber.Map{"data": "Product Header and Lines handled Successfully"}, fiber.StatusOK
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

func SaveDirectCommissionTransaction(distribId string, bvValue float64, reference string) (fiber.Map, int) {

	value := bvValue * 2.4

	refDistribId, err := repositories.GetRefDistribIdByDistribId(distribId)
	if err.Error != nil {
		return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError
	}

	obj := models.DirectCommissionTransaction{
		DistribId: refDistribId,
		Value:     value,
		Reference: reference,
	}

	if res := repositories.SaveDirectCommissionTransaction(obj); res.Error != nil {
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": "Direct Commission transaction successful"}, fiber.StatusOK
}

func handlePlaceOrderICoupons(AppliedCoupons []dto.PlaceOrderCoupon, distribId string, reference string, totalOrderAmount float64, tx *gorm.DB) (fiber.Map, int) {

	var (
		totalICouponBalance float64
	)

	//Get Total ICoupon Balance and save transaction loop
	for _, orderCoupon := range AppliedCoupons {
		balance, result := repositories.GetICouponBalance(orderCoupon.VID)
		if result.Error != nil {
			tx.Rollback()
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
		}

		if balance == 0 {
			repositories.CloseCoupon(orderCoupon.VID)
			tx.Rollback()
			return fiber.Map{
				"error": "ICoupon used already!",
			}, fiber.StatusPaymentRequired
		}

		if totalOrderAmount >= balance {

			totalOrderAmount = totalOrderAmount - balance

			// Updating in ICoupon Transaction table
			ICouponObj := models.ICouponTransaction{
				DistribId: distribId,
				VID:       orderCoupon.VID,
				Reference: reference,
				Value:     -balance, //using the full balance of ICoupon
			}

			if err := repositories.SaveICouponTx(ICouponObj); err.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError
			}

			//Balance will be zero after using full coupon, so active set to false
			if res := repositories.CloseCoupon(orderCoupon.VID); res.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
			}
		} else {

			ICouponObj := models.ICouponTransaction{
				DistribId: distribId,
				VID:       orderCoupon.VID,
				Value:     -totalOrderAmount, //The order amount is less than Icoupon Balance, so icoupon have remaining balance
			}
			if err := repositories.SaveICouponTx(ICouponObj); err.Error != nil {
				tx.Rollback()
				return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError
			}
			break //No need to loop again, since the order amount is satisfied with the coupon
		}

		totalICouponBalance += balance
	}

	if totalICouponBalance < totalOrderAmount {
		tx.Rollback()
		return fiber.Map{"error": "Insufficient Balance!"}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": "Place Coupons handled successfully"}, fiber.StatusOK
}

func GetOrderDetails(distrib_id string) (dto.OrderDetailsOut, int) {
	var subTotal float64 = 0.0
	var totalSandH float64 = 0.0
	var quantity uint = 0
	var orderProductArray []dto.OrderProduct
	var orderDetails dto.OrderDetailsOut
	var cartItems []dto.ProductsOut
	var userData models.User
	var TotalTypeValue float64
	//1. Retrieving All Products in Cart
	cartItems, result := repositories.GetAllCartProductsByDistribID(distrib_id)

	if result.Error != nil {
		return orderDetails, http.StatusInternalServerError
	}

	//2. Populating OrderProduct Array field
	for _, item := range cartItems {
		orderProduct := dto.OrderProduct{
			Name:        item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   uint64(item.Product.Price),
			SandH:       item.Product.SandH,
			SubTotal:    item.Product.Price * float64(item.Quantity),
			ProductType: item.Product.ProductType,
			TypeValue:   item.Product.TypeValue,
		}

		orderProductArray = append(orderProductArray, orderProduct)
		subTotal += orderProduct.SubTotal
		totalSandH += orderProduct.SandH
		quantity += item.Quantity
		TotalTypeValue += orderProduct.TypeValue * float64(item.Quantity)
	}
	//Retrieving User Data for Delivery Address
	userData, result = repositories.GetUserByID(distrib_id, userData)
	if result.Error != nil {
		println(fiber.Map{"error": result.Error})
		return orderDetails, http.StatusInternalServerError
	}

	deliveryAddress := dto.DeliveryAddress{
		ContactName:   userData.Name,
		ContactEmail:  userData.EmailAddress,
		Address:       userData.Address1,
		City:          userData.TownOrCity,
		District:      userData.District,
		State:         userData.StateOrProvince,
		ZipCode:       userData.PinOrZipCode,
		Country:       userData.Country,
		HomePhoneNo:   userData.HomePhoneNo,
		MobilePhoneNo: userData.MobilePhoneNo,
	}
	orderDetails = dto.OrderDetailsOut{
		Items:           orderProductArray,
		SubTotal:        subTotal,
		TotalSandH:      totalSandH,
		TotalAmount:     subTotal + totalSandH,
		DeliveryAddress: deliveryAddress,
		TotalQuantity:   float64(quantity),
		DistribId:       distrib_id,
		TotalTypeValue:  TotalTypeValue,
	}

	if result.Error != nil {
		print(fiber.Map{"error": result.Error})
		return orderDetails, http.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		print(fiber.Map{"data": "No Products in Cart"})
		return orderDetails, http.StatusNoContent
	}

	return orderDetails, http.StatusOK
}
