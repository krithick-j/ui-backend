package dto

type DeliveryAddress struct {
	ContactName   string
	ContactEmail  string
	Address       string
	City          string
	District      string
	State         string
	ZipCode       uint64
	Country       string
	HomePhoneNo   string
	MobilePhoneNo string
}

type OrderProduct struct {
	Name      string
	Quantity  uint
	UnitPrice uint64
	BV        int
	SubTotal  float64
	SandH     float64
}

type OrderDetailsOut struct {
	Items           []OrderProduct
	SubTotal        float64
	TotalSandH      float64
	TotalAmount     float64
	TotalQuantity   float64
	TotalBV         int
	DeliveryAddress DeliveryAddress
}

type PlaceOrderCoupon struct {
	VID string `json:"v_id"`
}

type ProductId struct {
	Productid uint `json:"product_id"`
}
type PlaceOrderIn struct {
	Products []ProductId        `json:"products"`
	Coupons  []PlaceOrderCoupon `json:"coupons"`
	// TotalAmount float64            `json:"total_amount"`
}
