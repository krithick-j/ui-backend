package dto

type AddTc struct {
	DistribID          string `json:"distrib_id"`
	PlacementDistribId string `json:"placement_distrib_id"`
	PlacementPlace     string `json:"placement_place"`
	PlacementSide      string `json:"placement_side"`
}
