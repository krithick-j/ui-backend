package dto

type AddTc struct {
	AppliedCoupons []PlaceOrderCoupon `json:"applied_coupons"`
	DistribID      string             `json:"distrib_id"`
	RefDistribId   string             `json:"ref_distrib_id"`
	RefPlace       string             `json:"ref_place"`
	RefSide        string             `json:"ref_side"`
}
