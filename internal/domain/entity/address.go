package domain

type Address struct {
	City         string  `json:"city"`
	State        string  `json:"state"`
	Neighborhood string  `json:"neighborhood"`
	Postcode     string  `json:"postcode"`
	CoordinateX  float32 `json:"coordinate_x"`
	CoordinateY  float32 `json:"coordinate_y"`
}
