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
