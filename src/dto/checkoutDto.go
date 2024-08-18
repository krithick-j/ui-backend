package dto

type Tc struct {
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	LPoint    int    `json:"l_point"`
	RPoint    int    `json:"r_point"`
	BvPoint   int    `json:"bv_point"`
}
type ChequeAvailableIn struct {
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	LTc       Tc     `json:"l_tc"`
	RTc       Tc     `json:"r_tc"`
}

type TakeChequeIn struct {
	DistribId   string `json:"distrib_id"`
	Place       string `json:"place"`
	ChequeCount int    `json:"cheque_count"`
}

type TakeChequeOut struct {
	AvailableChequeCountArr []ChequeFrequency `json:"cheque_count"`
}

type ChequeFrequency struct {
	Tc          string  `json:"tc"`
	Frequency   int     `json:"frequency"`
}

type PlacePointsArr struct {
	Place string  `json:"place,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type CheckoutIn struct {
	DistribId string `json:"distrib_id"`
}
type FrequencyForTc struct {
	DistribId string `json:"distrib_id"`
	Frequency int    `json:"frequency"`
}

type CheckoutFrequency struct {
	TotalCheckoutFrequency  int `json:"total_checkout_frequency"`
	ParentCheckoutFrequency int `json:"parent_checkout_frequency"`
	LeftCheckoutFrequency   int `json:"left_checkout_frequency"`
	RightCheckoutFrequency  int `json:"right_checkout_frequency"`
}

type ChequePinIn struct {
	DistribId  string `json:"distrib_id"`
	CurrentPin string `json:"current_pin"`
	NewPin     string `json:"new_pin"`
}

type ChequeLogin struct {
	DistribId string `json:"distrib_id"`
	Pin       string `json:"pin"`
}

type TcChequeFrequency struct {
	Place     string
	LPoint    int
	RPoint    int
	BvPoint   int
	Frequency uint
}
