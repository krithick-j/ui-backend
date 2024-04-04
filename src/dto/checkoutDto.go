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
	DistribId string `json:"distrib_id"`
	Place string `json:"place"`
}

type CheckoutIn struct {
	DistribId string `json:"distrib_id"`
}