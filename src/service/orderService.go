package service

import (
	"net/http"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PlaceOrder(OrderIn dto.PlaceOrderIn, tx *gorm.DB) (fiber.Map, int) {

	orderId := GenerateUniqueHexCode(10)
	configs.Log.Infof("Generating a new order id %s", orderId)
	//sum product value
	res, status := GetOrderDetails(OrderIn.DistribId, tx)
	total := res["data"].(dto.OrderDetailsOut)
	total.OrderId = orderId
	if status != fiber.StatusOK {
		return utils.NotNilErrorMessage(res["error"].(error), "GetOrderDetails", "PlaceOrder", status, tx)
	}

	totalOrderAmount := total.TotalAmount

	productType := total.Products[0].ProductType
	configs.Log.Infof("Product type %v", productType)
	if res, status := handleProductHeaderAndLines(OrderIn, orderId, total, productType, tx); status != fiber.StatusOK {
		return utils.NotNilErrorMessage(res["error"].(error), "handleProductHeaderAndLines", "PlaceOrder", status, tx)
	}

	if res, status := handleProductType(OrderIn, productType, orderId, totalOrderAmount, tx, total); status != fiber.StatusOK {
		return utils.NotNilErrorMessage(res["error"].(error), "handleProductType", "PlaceOrder", status, tx)
	}

	invoicePdfPath, status, err := GenerateInvoice(orderId, productType, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GenerateInvoice", "PlaceOrder", status, tx)
	}

	err = SendHtmlMailOrder(total, invoicePdfPath)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SendHtmlMailOrder", "PlaceOrder", fiber.StatusInternalServerError, tx)
	}
	configs.Log.Infoln("Mail Sent!")
	return fiber.Map{"success": "Ordered Placed Successfully", "link": invoicePdfPath}, fiber.StatusOK
}

func handleProductType(OrderIn dto.PlaceOrderIn, productType string, orderId string, totalOrderAmount float64, tx *gorm.DB, total dto.OrderDetailsOut) (fiber.Map, int) {
	configs.Log.Infoln("Starting BV Product")
	if OrderIn.AppliedCoupons != nil && productType != "ep" {

		if res, status := handlePlaceOrderICoupons(OrderIn.AppliedCoupons, OrderIn.DistribId, orderId, totalOrderAmount, tx); status != fiber.StatusOK {
			configs.Log.Errorln("Error on handlePlaceOrderICoupons from handleProductType fn")
			return fiber.Map{"error": res["error"]}, status
		}

		configs.Log.Infoln("Coupon applied successfully")
		configs.Log.Infof("The product type is %s", productType)
		if productType == "bv" {
			configs.Log.Infoln("The product is of bv type")
			if res, status := handleBvProduct(OrderIn, orderId, total, tx); status != fiber.StatusOK {
				configs.Log.Errorln("Error on handleBvProduct from handleProductType fn")
				return fiber.Map{"error": res["error"]}, fiber.StatusInternalServerError
			}

		} else if productType == "rsp" {
			configs.Log.Infoln("Handling RSP value-->")
			configs.Log.Infoln("rsp value-->", total.TotalTypeValue)
			if err := repositories.AddRsp(OrderIn.DistribId, orderId, total.TotalTypeValue, tx); err != nil {
				tx.Rollback()
				configs.Log.Errorln("Error on AddRspTx from handleProductType fn")
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}
		} else {
			tx.Rollback()
			configs.Log.Errorln("Error on product Type from handleProductType fn")
			return fiber.Map{"error": "Product type does not match"}, fiber.StatusInternalServerError
		}

	} else if productType == "ep" {

		if res, status := handleEpProduct(total, OrderIn, orderId, tx); status != fiber.StatusOK {
			configs.Log.Errorln("Error on handleEpProduct from handleProductType fn")
			return fiber.Map{"error": res["error"]}, fiber.StatusInternalServerError
		}
	}

	_, err := repositories.DeleteAllCartProduct(OrderIn.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on DeleteAllCartProduct repositories from handleProductType service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	return nil, fiber.StatusOK
}

func handleEpProduct(total dto.OrderDetailsOut, OrderIn dto.PlaceOrderIn, orderId string, tx *gorm.DB) (fiber.Map, int) {
	totalEp, err := repositories.GetEpBalance(total.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on GetEpBalance from handleEpProduct fn")
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	if totalEp >= total.TotalTypeValue {
		err := repositories.SaveEpTx(OrderIn.DistribId, orderId, total.TotalTypeValue, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on SaveEpTx from handleEpProduct fn: ", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
	} else {
		tx.Rollback()
		configs.Log.Warnln("Insufficient Balance!")
		return fiber.Map{"data": "Insufficient Balance!"}, fiber.StatusForbidden
	}

	return fiber.Map{"data": "Handling Ep product successful"}, fiber.StatusOK
}

func handleBvProduct(OrderIn dto.PlaceOrderIn, orderId string, total dto.OrderDetailsOut, tx *gorm.DB) (fiber.Map, int) {
	configs.Log.Infoln("Preparing BV product handling...")

	//Referral distrib id is taken
	referralDistribId, err := repositories.GetRefDistribIdByDistribId(OrderIn.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId repositories fn from handleBvProduct service fn")
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	err = repositories.AddDirectBvTx(referralDistribId, orderId, total.TotalTypeValue, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "AddDirectBvTx", "handleBvProduct", fiber.StatusInternalServerError, tx)
	}
	configs.Log.Infoln("Direct BV Added")

	msg, status := UpdateCurrentPlaceValues(OrderIn.DistribId, OrderIn.PlaceBvs, orderId, tx, total.TotalTypeValue)
	if status == fiber.StatusInternalServerError {
		return fiber.Map{"error": msg}, status
	}
	configs.Log.Infoln("Place BV Added")
	configs.Log.Infoln("Direct Commission transaction")

	Dcmessage, status := SaveDirectCommissionTransaction(OrderIn.DistribId, total.TotalTypeValue, orderId, tx)
	if status != fiber.StatusOK {
		return fiber.Map{"error": Dcmessage["error"]}, status
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
	err := repositories.SaveOrderHeader(OrderHeaderObj, tx)
	configs.Log.Infoln("OrdersHeaders: ", OrderHeaderObj.ID)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SaveOrderHeader repositories fn from handleProductHeaderAndLines service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	configs.Log.Infoln("Committing transaction of order headers")

	for _, placeBv := range OrderIn.PlaceBvs {
		addBv := models.AddedBv{
			OrderHeaderID: OrderHeaderObj.ID,
			ReferenceNo:   orderId,
			DistribId:     OrderIn.DistribId,
			Place:         placeBv.Place,
			Value:         placeBv.AddBv,
		}
		err = repositories.SaveAddedbv(&addBv, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling SaveAddedbv repositories fn from handleProductHeaderAndLines service fn: ", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
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

		configs.Log.Infoln("Orders Liner Obj: \n", OrderLinerObj)
		err = repositories.SaveOrderLiner(OrderLinerObj, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling SaveOrderLiner repositories fn from handleProductHeaderAndLines service fn: ", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
	}
	return fiber.Map{"data": "Product Header and Lines handled Successfully"}, fiber.StatusOK
}

func GetAllOrders(tx *gorm.DB) (fiber.Map, int) {

	var OrdersOut []dto.AllOrdersOut

	configs.Log.Infoln("Retrieving All Orders INIT")
	order, err := repositories.GetAllOrders(tx)

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Warn("No Orders Found")
		return fiber.Map{"data": "Not Found"}, http.StatusNotFound
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

		for _, productLine := range order.OrdersLiners {
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

	return fiber.Map{"data": OrdersOut}, fiber.StatusOK
}

func GetOrdersByDistribId(distribId string, tx *gorm.DB) (fiber.Map, int) {

	var OrdersOut []dto.AllOrdersOut

	configs.Log.Infoln("Retrieving AllOrdersByDistribId INIT")
	orders, err := repositories.GetOrderByDistribId(distribId, tx)
	if err == gorm.ErrRecordNotFound {
		configs.Log.Warnln("No Orders Found on calling GetOrderByDistribId: ", err.Error())
		return fiber.Map{"data": "Not Found", "err": err}, fiber.StatusNotFound
	}
	if err != nil {
		configs.Log.Errorf("%v", err.Error())
		return fiber.Map{"error": err}, fiber.StatusInternalServerError
	}
	configs.Log.Infoln("Retrieving All Orders DONE")
	configs.Log.Infoln("Retrieving Orders: \n", orders)

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

		for _, productLine := range order.OrdersLiners {
			imageUrl, err := repositories.GetProductImage(productLine.ProductID, tx)
			if err != nil && err != gorm.ErrRecordNotFound {
				return utils.NotNilErrorMessage(err, "GetProductImage", "GetOrdersByDistribId", fiber.StatusInternalServerError, tx)
			}
			productArrObj := dto.ProductDetails{
				Id:           productLine.ID,
				Name:         productLine.Name,
				Quantity:     productLine.Quantity,
				Price:        productLine.UnitPrice,
				ProductImage: imageUrl,
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
	return fiber.Map{"data": OrdersOut}, fiber.StatusOK
}

func SaveDirectCommissionTransaction(distribId string, bvValue float64, reference string, tx *gorm.DB) (fiber.Map, int) {

	value := bvValue * 2.4

	refDistribId, err := repositories.GetRefDistribIdByDistribId(distribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId repositories fn from SaveDirectCommissionTransaction service fn")
		return fiber.Map{"error": err.Error(), "err": err}, fiber.StatusInternalServerError
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

	if err := repositories.SaveDirectCommissionTransaction(obj, tx); err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRefDistribIdByDistribId repositories fn from SaveDirectCommissionTransaction service fn")
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": "Direct Commission transaction successful"}, fiber.StatusOK
}

func handlePlaceOrderICoupons(AppliedCoupons []dto.PlaceOrderCoupon, distribId string, reference string, totalOrderAmount float64, tx *gorm.DB) (fiber.Map, int) {

	var (
		totalICouponBalance float64
	)

	configs.Log.Infoln("Checking the ICoupons")
	for _, orderCoupon := range AppliedCoupons {
		expiresOn, err := repositories.GetICouponExpiryDate(orderCoupon.VID, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error Retrieving the ICoupon Expiry Date")
			return fiber.Map{"error": err.Error()}, fiber.StatusBadRequest
		}

		if expiresOn.Before(time.Now()) || expiresOn.Equal(time.Now()) {
			tx.Rollback()
			configs.Log.Errorln("ICoupon expired")
			return fiber.Map{"error": "ICoupon expired"}, fiber.StatusBadRequest
		}

		balance, err := repositories.GetICouponBalance(orderCoupon.VID, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling GetICouponBalance repositories fn from handlePlaceOrderICoupons service fn")
			return fiber.Map{"error": err.Error()}, fiber.StatusBadRequest
		}

		if balance == 0 {
			repositories.CloseCoupon(orderCoupon.VID, tx)
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

			if err := repositories.SaveICouponTx(ICouponObj, tx); err != nil {
				return utils.NotNilErrorMessage(err, "SaveICouponTx", "handlePlaceOrderICoupons", fiber.StatusBadRequest, tx)
			}

			//Balance will be zero after using full coupon, so active set to false
			if err := repositories.CloseCoupon(orderCoupon.VID, tx); err != nil {
				return utils.NotNilErrorMessage(err, "CloseCoupon", "handlePlaceOrderICoupons", fiber.StatusBadRequest, tx)
			}
		} else {

			ICouponObj := models.ICouponTransaction{
				DistribId: distribId,
				VID:       orderCoupon.VID,
				Value:     -totalOrderAmount, //The order amount is less than Icoupon Balance, so icoupon have remaining balance
				Reference: reference,
			}
			if err := repositories.SaveICouponTx(ICouponObj, tx); err != nil {
				return utils.NotNilErrorMessage(err, "SaveICouponTx", "handlePlaceOrderICoupons", fiber.StatusBadRequest, tx)
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

// dto.OrderDetailsOut
func GetOrderDetails(distrib_id string, tx *gorm.DB) (fiber.Map, int) {
	var subTotal float64 = 0.0
	var totalSandH float64 = 0.0
	var quantity uint = 0
	var orderProductArray []dto.OrderProduct
	var orderDetails dto.OrderDetailsOut
	var cartItems []dto.ProductsOut
	var TotalTypeValue float64

	//1. Retrieving All Products in Cart
	cartItems, err := repositories.GetAllCartProductsByDistribID(distrib_id, tx)
	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllCartProductsByDistribID", "GetOrderDetails", fiber.StatusBadRequest, tx)
	}
	configs.Log.Infof("Cart items --> %v", cartItems)

	//2. Populating OrderProduct Array field
	for _, item := range cartItems {
		imageUrl, err := repositories.GetProductImage(item.Product.ID, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetProductImage", "GetOrderDetails", fiber.StatusBadRequest, tx)

		}
		configs.Log.Infof("The Individual Item %v", item.Product.ProductType)
		orderProduct := dto.OrderProduct{
			ProductID:     item.Product.ID,
			ProductImage:  imageUrl,
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
	userData, err := repositories.GetUserByID(distrib_id, tx)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GetOrderDetails fn: ", err.Error())
		return fiber.Map{"error": err.Error(), "err": err}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetUserByID repositories fn from GetOrderDetails fn: ", err.Error())
		return fiber.Map{"error": err.Error(), "err": err}, fiber.StatusInternalServerError
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

	configs.Log.Infof("%v", orderDetails)
	return utils.SuccessMessage(orderDetails, fiber.StatusOK)
}
