package dto

import "ui-back-end/src/models"

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
	Name      string  `json:"name"`
	Quantity  uint    `json:"quantity"`
	UnitPrice uint64  `json:"unit_price"`
	BV        int     `json:"bv"`
	SubTotal  float64 `json:"sub_total"`
	SandH     float64 `json:"s_and_h"`
	Rsp       int     `json:"rsp"`
	Ep        float64 `json:"ep"`
}

type OrderDetailsOut struct {
	DistribId       string
	Items           []OrderProduct  `json:"items"`
	SubTotal        float64         `json:"sub_total"`
	TotalSandH      float64         `json:"total_s_and_h"`
	TotalAmount     float64         `json:"total_amount"`
	TotalQuantity   float64         `json:"total_quantity"`
	TotalBV         int             `json:"total_bv"`
	TotalRsp        int             `json:"total_rsp"`
	TotalEp         float64         `json:"total_ep"`
	DeliveryAddress DeliveryAddress `json:"delivery_address"`
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
	Place string `json:"place"`
	AddBv int    `json:"add_bv"`
}
type PlaceOrderIn struct {
	AppliedCoupons []PlaceOrderCoupon `json:"applied_coupons"` //Coupons Applied
	DistribId      string             `json:"distrib_id"`
	TotalAmount    float64            `json:"total_amount"`
	PlaceBvs       []PlaceBv          `json:"place_bvs"`
	OrderType      string             `json:"order_type"`
}

type OrdersOut struct {
	OrderHeader models.OrdersHeader
	OrderLiners models.OrdersLiner
}
