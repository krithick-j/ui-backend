package service

import (
	"fmt"
	"net/http"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn) (fiber.Map, int) {

	orderId := GenerateUniqueHexCode(10)
	configs.Log.Infof("Generating a new order id %s", orderId)

	//sum product value
	total, status := GetOrderDetails(OrderIn.DistribId)
	total.OrderId = orderId
	if status != http.StatusOK {
		return fiber.Map{"data": "Something gone wrong"}, fiber.StatusInternalServerError
	}

	tx := configs.DB.Begin()

	totalOrderAmount := total.TotalAmount

	productType := total.Products[0].ProductType
	configs.Log.Infof("Product type %v", productType)
	if res, status := handleProductHeaderAndLines(OrderIn, orderId, total, productType, tx); status != fiber.StatusOK {
		return fiber.Map{"error": res["error"]}, status
	}

	if res, status := handleProductType(OrderIn, productType, orderId, totalOrderAmount, tx, total); status != fiber.StatusOK {
		return fiber.Map{"error": res["error"]}, status
	}

	if commitRes := tx.Commit(); commitRes.Error != nil {
		return fiber.Map{"error": commitRes.Error.Error()}, fiber.StatusInternalServerError
	}

	//SendHtmlMailOrder(dto.OrderDetailsOut{})
	return fiber.Map{"success": "Ordered Placed Successfully"}, http.StatusOK
}

func handleProductType(OrderIn dto.PlaceOrderIn, productType string, orderId string, totalOrderAmount float64, tx *gorm.DB, total dto.OrderDetailsOut) (fiber.Map, int) {
	configs.Log.Infoln("Starting BV Product")
	if OrderIn.AppliedCoupons != nil && productType != "ep" {

		if res, status := handlePlaceOrderICoupons(OrderIn.AppliedCoupons, OrderIn.DistribId, orderId, totalOrderAmount, tx); status != fiber.StatusOK {
			return fiber.Map{"error": res["error"]}, status
		}

		configs.Log.Infoln("Coupon applied successfully")
		configs.Log.Infof("The product type is %s", productType)
		if productType == "bv" {
			configs.Log.Infoln("The product is of bv type")
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
	invoicePdfPath, status, err := GenerateInvoice(orderId)
	if err != nil {
		return fiber.Map{"error": err.Error()}, status
	}

	err = SendHtmlMailOrder(total, invoicePdfPath)
	if err != nil {
		configs.Log.Errorf("Error sending email : %s", err.Error())
	}
	configs.Log.Infoln("Mail Sent!")
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
	configs.Log.Infoln("Preparing BV product handling...")

	// res := repositories.AddDirectBvTx(OrderIn.DistribId, orderId, total.TotalTypeValue)
	// if res.Error != nil {
	// 	tx.Rollback()
	// 	return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	// }
	// configs.Log.Infoln("Direct BV Added")

	msg, status := UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId, tx, total.TotalTypeValue)
	if status == fiber.StatusInternalServerError {
		return fiber.Map{"error": msg}, status
	}
	configs.Log.Infoln("Place BV Added")
	configs.Log.Infoln("Direct Commission transaction")

	Dcmessage, status := SaveDirectCommissionTransaction(OrderIn.DistribId, total.TotalTypeValue, orderId)
	if status != 200 {
		return fiber.Map{"error": Dcmessage}, status
	}
	configs.Log.Infoln("Direct Commison Added")
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
		ContactName:    total.CustomerDetails.Name,
		ContactEmail:   total.CustomerDetails.Email,
		Address:        total.ShippingAddress.Address,
		City:           total.ShippingAddress.City,
		District:       total.ShippingAddress.District,
		State:          total.ShippingAddress.State,
		ZipCode:        total.ShippingAddress.ZipCode,
		Country:        total.ShippingAddress.Country,
		HomePhoneNo:    total.CustomerDetails.HomePhoneNo,
		MobilePhoneNo:  total.CustomerDetails.MobilePhoneNo,
	}
	err := repositories.SaveOrderHeader(OrderHeaderObj)
	if err != nil {
		tx.Rollback()
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	for _, placeBv := range OrderIn.PlaceBvs {
		addBv := models.AddedBv{
			OrderHeaderID: OrderHeaderObj.ID,
			ReferenceNo:   orderId,
			DistribId:     OrderIn.DistribId,
			Place:         placeBv.Place,
			Value:         placeBv.AddBv,
		}
		repositories.SaveAddedbv(&addBv)
	}

	configs.Log.Infoln("Products", total.Products)
	for _, product := range total.Products {
		configs.Log.Infoln("Product Line: ", product)
		OrderLinerObj := &models.OrdersLiner{
			OrdersHeaderID: OrderHeaderObj.ID,
			ProductID:      product.ProductID,
			Name:           product.Name,
			Quantity:       product.Quantity,
			UnitPrice:      product.UnitPrice,
			ProductType:    product.ProductType,
			TypeValue:      product.TypeValue,
			SubTotal:       product.SubTotal,
			SandH:          product.SandH,
			GstPercentage:  product.GstPercentage,
		}

		err = repositories.SaveOrderLiner(OrderLinerObj)
		if err != nil {
			tx.Rollback()
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
	}
	return fiber.Map{"data": "Product Header and Lines handled Successfully"}, fiber.StatusOK
}

func GetAllOrders() (fiber.Map, int) {

	var OrdersOut []dto.AllOrdersOut

	configs.Log.Infoln("Retrieving All Orders INIT")
	order, result := repositories.GetAllOrders()

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("No Orders Found")
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorf("%v", result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, http.StatusInternalServerError
	}
	configs.Log.Infoln("Retrieving All Orders DONE")
	for _, order := range order {
		var productArr []dto.ProductDetails

		ShippingAddressObj := dto.ShippingAddress{
			Address:  order.Address,
			City:     order.City,
			District: order.District,
			State:    order.State,
			ZipCode:  order.ZipCode,
			Country:  order.Country,
		}

		CustomerDetailsObj := dto.CustomerDetails{
			DistribId:     order.DistribId,
			Name:          order.ContactName,
			Email:         order.ContactEmail,
			MobilePhoneNo: order.MobilePhoneNo,
			HomePhoneNo:   order.HomePhoneNo,
		}

		for _, productLine := range order.OrdersLiner {
			productArrObj := dto.ProductDetails{
				Id:           productLine.ID,
				Name:         productLine.Name,
				Quantity:     productLine.Quantity,
				Price:        productLine.UnitPrice,
				ProductImage: "",
				ProductType:  productLine.ProductType,
				SAndH:        productLine.SandH,
				SubTotal:     productLine.SubTotal,
				TypeValue:    productLine.TypeValue,
			}
			productArr = append(productArr, productArrObj)
		}

		OrderArr := dto.AllOrdersOut{
			OrderId:            order.OrderId,
			SubTotal:           order.SubTotal,
			TotalAmount:        order.TotalAmount,
			TotalSandH:         order.TotalSandH,
			DeliveryStatus:     order.DeliveryStatus,
			ShipmentTrackingNo: order.ShipmentTrackingNo,
			CourierName:        order.CourierName,
			CreatedAt:          order.CreatedAt.UTC().String(),
			UpdatedAt:          order.UpdatedAt.UTC().String(),
			DeletedAt:          order.DeletedAt.Time.UTC().String(),
			DeliveredAt:        order.DeliveredAt,
			ShippingAddress:    ShippingAddressObj,
			CustomerDetails:    CustomerDetailsObj,
			ProductDetails:     productArr,
		}

		OrdersOut = append(OrdersOut, OrderArr)
	}
	return fiber.Map{"data": OrdersOut}, http.StatusOK
}

func GetOrdersByDistribId(distribId string) (fiber.Map, int) {

	var OrdersOut []dto.AllOrdersOut

	configs.Log.Infoln("Retrieving AllOrdersByDistribId INIT")
	orders, result := repositories.GetOrderByDistribId(distribId)

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("No Orders Found")
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
	}

	if result.Error != nil {
		configs.Log.Errorf("%v", result.Error.Error())
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	configs.Log.Infoln("Retrieving All Orders DONE")
	for _, order := range orders {
		var productArr []dto.ProductDetails

		ShippingAddressObj := dto.ShippingAddress{
			Address:  order.Address,
			City:     order.City,
			District: order.District,
			State:    order.State,
			ZipCode:  order.ZipCode,
			Country:  order.Country,
		}

		CustomerDetailsObj := dto.CustomerDetails{
			DistribId:     order.DistribId,
			Name:          order.ContactName,
			Email:         order.ContactEmail,
			MobilePhoneNo: order.MobilePhoneNo,
			HomePhoneNo:   order.HomePhoneNo,
		}

		for _, productLine := range order.OrdersLiner {
			productArrObj := dto.ProductDetails{
				Id:           productLine.ID,
				Name:         productLine.Name,
				Quantity:     productLine.Quantity,
				Price:        productLine.UnitPrice,
				ProductImage: "",
				ProductType:  productLine.ProductType,
				SAndH:        productLine.SandH,
				SubTotal:     productLine.SubTotal,
				TypeValue:    productLine.TypeValue,
			}
			productArr = append(productArr, productArrObj)
		}

		OrderArr := dto.AllOrdersOut{
			OrderId:            order.OrderId,
			SubTotal:           order.SubTotal,
			TotalAmount:        order.TotalAmount,
			TotalSandH:         order.TotalSandH,
			DeliveryStatus:     order.DeliveryStatus,
			ShipmentTrackingNo: order.ShipmentTrackingNo,
			CourierName:        order.CourierName,
			CreatedAt:          order.CreatedAt.UTC().String(),
			UpdatedAt:          order.UpdatedAt.UTC().String(),
			DeletedAt:          order.DeletedAt.Time.UTC().String(),
			DeliveredAt:        order.DeliveredAt,
			TotalQuantity:      order.TotalQuantity,
			TotalTypeValue:     order.TotalTypeValue,
			ShippingAddress:    ShippingAddressObj,
			CustomerDetails:    CustomerDetailsObj,
			ProductDetails:     productArr,
		}

		OrdersOut = append(OrdersOut, OrderArr)
	}
	return fiber.Map{"data": OrdersOut}, http.StatusOK
}

func SaveDirectCommissionTransaction(distribId string, bvValue float64, reference string) (fiber.Map, int) {

	value := bvValue * 2.4

	refDistribId, err := repositories.GetRefDistribIdByDistribId(distribId)
	fmt.Println("reference distrib id ", refDistribId)
	if err.Error != nil {
		return fiber.Map{"error": err.Error.Error()}, fiber.StatusInternalServerError

	}
	activateDayNumber := 21
	obj := models.DirectCommissionTransaction{
		DistribId:    distribId,
		Value:        value,
		Reference:    reference,
		RefDistribId: refDistribId,
		ActivateDate: time.Now().AddDate(0, 0, activateDayNumber), //21 days
		ExpiryDate:   time.Now().AddDate(0, 6, activateDayNumber), //6 months
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

	configs.Log.Infoln("Checking the ICoupons")
	for _, orderCoupon := range AppliedCoupons {
		expiresOn, result := repositories.GetICouponExpiryDate(orderCoupon.VID)
		if result.Error != nil {
			tx.Rollback()
			configs.Log.Errorln("Error Retrieving the ICoupon Expiry Date")
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusBadRequest
		}

		if expiresOn.Before(time.Now()) || expiresOn.Equal(time.Now()) {
			tx.Rollback()
			configs.Log.Errorln("ICoupon Expired")
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusBadRequest
		}

		balance, result := repositories.GetICouponBalance(orderCoupon.VID)
		if result.Error != nil {
			tx.Rollback()
			configs.Log.Errorln("Error Retrieving the ICoupon Balance")
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusBadRequest
		}

		if balance == 0 {
			repositories.CloseCoupon(orderCoupon.VID)
			tx.Rollback()
			configs.Log.Errorln("Error Detecting used Icoupons")
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
				configs.Log.Errorln("Error while saving Icoupon Balance Transaction")

				return fiber.Map{"error": err.Error.Error()}, fiber.StatusBadRequest
			}

			//Balance will be zero after using full coupon, so active set to false
			if res := repositories.CloseCoupon(orderCoupon.VID); res.Error != nil {
				tx.Rollback()
				configs.Log.Errorf("Error while closing the icoupon, %s", res.Error.Error())

				return fiber.Map{"error": res.Error.Error()}, fiber.StatusBadRequest
			}
		} else {

			ICouponObj := models.ICouponTransaction{
				DistribId: distribId,
				VID:       orderCoupon.VID,
				Value:     -totalOrderAmount, //The order amount is less than Icoupon Balance, so icoupon have remaining balance
				Reference: reference,
			}
			if err := repositories.SaveICouponTx(ICouponObj); err.Error != nil {
				tx.Rollback()
				configs.Log.Errorw("Error while saving icoupon transactions %v", ICouponObj)
				return fiber.Map{"error": err.Error.Error()}, fiber.StatusBadRequest
			}
			totalICouponBalance += balance
			configs.Log.Infoln("TotalICouponBalance +=", totalICouponBalance)
			break //No need to loop again, since the order amount is satisfied with the coupon
		}
		configs.Log.Infoln("ICoupon ", orderCoupon.VID, " checking done")

		totalICouponBalance += balance
		configs.Log.Infoln("TotalICouponBalance +=", totalICouponBalance)
	}
	configs.Log.Infoln("Checking ICoupon values done")
	configs.Log.Infoln("Checking ICoupon value with total value")
	configs.Log.Infoln("totalICoupon Balance-> ", totalICouponBalance, " totalOrderAmount-> ", totalOrderAmount)

	if totalICouponBalance < totalOrderAmount {
		tx.Rollback()
		configs.Log.Errorln("Insufficient Balance")
		return fiber.Map{"error": "Insufficient Balance!"}, fiber.StatusInternalServerError
	}
	configs.Log.Infoln("Place Coupons handled successfully")
	return fiber.Map{"data": "Place Coupons handled successfully"}, fiber.StatusOK
}

func GetOrderDetails(distrib_id string) (dto.OrderDetailsOut, int) {
	var subTotal float64 = 0.0
	var totalSandH float64 = 0.0
	var quantity uint = 0
	var orderProductArray []dto.OrderProduct
	var orderDetails dto.OrderDetailsOut
	var cartItems []dto.ProductsOut
	var TotalTypeValue float64
	//1. Retrieving All Products in Cart
	cartItems, result := repositories.GetAllCartProductsByDistribID(distrib_id)

	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetAllCartProductsByDistribID repositories fn from GetOrderDetails fn ", result.Error.Error())
		return orderDetails, fiber.StatusInternalServerError
	}
	configs.Log.Infof("%v", orderDetails.Products)
	if len(cartItems) < 1 {
		configs.Log.Warnln("No Items found on the card")
		return orderDetails, fiber.StatusNotFound
	}

	//2. Populating OrderProduct Array field
	for _, item := range cartItems {

		configs.Log.Infof("The Individual Item %v", item.Product.ProductType)
		orderProduct := dto.OrderProduct{
			ProductID:     item.Product.ID,
			ProductImage:  "",
			Name:          item.Product.Name,
			Quantity:      item.Quantity,
			UnitPrice:     item.Product.Price,
			SubTotal:      item.Product.Price * float64(item.Quantity),
			SandH:         item.Product.SandH,
			ProductType:   item.Product.ProductType,
			TypeValue:     item.Product.TypeValue,
			GstPercentage: item.Product.GstPercentage,
		}

		orderProductArray = append(orderProductArray, orderProduct)
		subTotal += orderProduct.SubTotal
		totalSandH += orderProduct.SandH * float64(item.Quantity)
		quantity += item.Quantity
		TotalTypeValue += orderProduct.TypeValue * float64(item.Quantity)
	}
	//Retrieving User Data for Delivery Address
	userData, result := repositories.GetUserByID(distrib_id)
	if result.Error != nil {
		configs.Log.Warnf("%v", orderDetails)
		return orderDetails, fiber.StatusBadRequest
	}

	shippingAddress := dto.ShippingAddress{
		Address:  userData.Address1,
		City:     userData.TownOrCity,
		District: userData.District,
		State:    userData.StateOrProvince,
		ZipCode:  userData.PinOrZipCode,
		Country:  userData.Country,
	}

	customerDetails := dto.CustomerDetails{
		DistribId:     distrib_id,
		Name:          userData.Name,
		Email:         userData.EmailAddress,
		MobilePhoneNo: userData.MobilePhoneNo,
		HomePhoneNo:   userData.HomePhoneNo,
	}

	orderDetails = dto.OrderDetailsOut{
		DistribId:       distrib_id,
		Products:        orderProductArray,
		SubTotal:        subTotal,
		TotalSandH:      totalSandH,
		TotalAmount:     subTotal + totalSandH,
		TotalQuantity:   float64(quantity),
		ShippingAddress: shippingAddress,
		CustomerDetails: customerDetails,
		TotalTypeValue:  TotalTypeValue,
	}

	if result.Error != nil {
		return orderDetails, fiber.StatusBadRequest
	}

	if result.RowsAffected == 0 {
		return orderDetails, http.StatusNoContent
	}
	configs.Log.Infof("%v", orderDetails)
	return orderDetails, http.StatusOK
}
