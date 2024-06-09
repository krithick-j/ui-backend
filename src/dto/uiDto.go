package dto

type BvCounterIn struct {
	DistribId string `json:"distrib_id"`
	Place     string `json:"place"`
	StartDate string `json:"start_date"`
}

type BvCounterOut struct {
	RankBv       LeftAndRight `json:"rank_bv"`
	CommissionBv LeftAndRight `json:"commission_bv"`
}

type LeftAndRight struct {
	Left  float64 `json:"left"`
	Right float64 `json:"right"`
}
