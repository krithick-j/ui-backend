package dto

import "time"

type Coupon struct {
	VID        string
	Value      uint
	ExpiryDate time.Time
}

// ICoupon Input from Admin
type ICouponIn struct {
	DistribID string
	DateOn    time.Time
	TxDetail  string
	AdminName string
	Coupons   []Coupon
}
