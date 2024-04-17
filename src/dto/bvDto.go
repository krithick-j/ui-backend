package dto

type TrackingCenterOut struct {
	DistribId string
	Tc        []TCBv
}

type TCBv struct {
    Place     string `json:"place"`
    LPoint    int    `json:"l_point"`
    RPoint    int    `json:"r_point"`
    BvPoint   int    `json:"bv_point"`
}
