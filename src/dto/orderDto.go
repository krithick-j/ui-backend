package dto

import (
	"ui-back-end/src/models"
)

type DeliveryAddress struct {
	ContactName   string `json:"contact_name"`
	ContactEmail  string `json:"contact_email"`
	Address       string `json:"address"`
	City          string `json:"city"`
	District      string `json:"district"`
	State         string `json:"state"`
	ZipCode       uint64 `json:"zip_code"`
	Country       string `json:"country"`
	HomePhoneNo   string `json:"home_phone_no"`
	MobilePhoneNo string `json:"mobile_phone_no"`
}

type OrderProduct struct {
	ProductID     uint    `json:"product_id"`
	ProductImage  string  `json:"product_image"`
	Name          string  `json:"name"`
	Quantity      uint    `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	SubTotal      float64 `json:"sub_total"`
	SandH         float64 `json:"s_and_h"`
	ProductType   string  `json:"product_type"`
	TypeValue     float64 `json:"type_value"`
	GstPercentage float64 `json:"gst_percentage"`
}

type OrderDetailsOut struct {
	DistribId       string
	Products        []OrderProduct  `json:"products"`
	SubTotal        float64         `json:"sub_total"`
	TotalSandH      float64         `json:"total_s_and_h"`
	TotalAmount     float64         `json:"total_amount"`
	TotalQuantity   float64         `json:"total_quantity"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
	CustomerDetails CustomerDetails `json:"customer_details"`
	TotalTypeValue  float64         `json:"total_type_value"`
	CreatedAt       string          `json:"created_at"`
	OrderId         string          `json:"order_id"`
}

type PlaceOrderCoupon struct {
	VID            string  `json:"v_id"`
	Pin            string  `json:"pin"`
	AmountDetected float64 `json:"amount_detected"`
}

type ProductId struct {
	Productid uint `json:"product_id"`
}

type PlaceBv struct {
	Place string  `json:"place"`
	AddBv float64 `json:"add_bv"`
}
type PlaceOrderIn struct {
	AppliedCoupons []PlaceOrderCoupon `json:"applied_coupons"` //Coupons Applied
	DistribId      string             `json:"distrib_id"`
	PlaceBvs       []PlaceBv          `json:"place_bvs"`
}

type OrdersOut struct {
	OrderHeader models.OrdersHeader
	OrderLiners models.OrdersLiner
}

type AllOrdersOut struct {
	OrderId            string           `json:"order_id"`
	SubTotal           float64          `json:"sub_total"`
	TotalAmount        float64          `json:"total_amount"`
	TotalSandH         float64          `json:"total_s_and_h"`
	DeliveryStatus     string           `json:"delivery_status"`
	ShipmentTrackingNo string           `json:"shipment_tracking_no"`
	CourierName        string           `json:"courier_name"`
	CreatedAt          string           `json:"created_at"`
	UpdatedAt          string           `json:"updated_at"`
	DeletedAt          string           `json:"deleted_at"`
	DeliveredAt        string           `json:"delivered_at"`
	ShippingAddress    ShippingAddress  `json:"shipping_address"`
	CustomerDetails    CustomerDetails  `json:"customer_details"`
	ProductDetails     []ProductDetails `json:"products"`
	TotalQuantity      float64          `json:"total_quantity"`
	TotalTypeValue     float64          `json:"total_type_value"`
}

type ShippingAddress struct {
	Address  string `json:"address"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
}

type CustomerDetails struct {
	DistribId     string `json:"distrib_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	MobilePhoneNo string `json:"mobile_phone_no"`
	HomePhoneNo   string `json:"home_phone_no"`
}

type ProductDetails struct {
	Id           uint    `json:"product_id"`
	Name         string  `json:"name"`
	Quantity     uint    `json:"quantity"`
	Price        float64 `json:"unit_price"`
	ProductImage string  `json:"product_image"`
	ProductType  string  `json:"product_type"`
	SAndH        float64 `json:"s_and_h"`
	SubTotal     float64 `json:"sub_total"`
	TypeValue    float64 `json:"type_value"`
}
