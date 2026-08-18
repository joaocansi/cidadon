package address

import "time"

type Address struct {
	ID           string    `json:"id"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Neighborhood string    `json:"neighborhood"`
	Postcode     string    `json:"postcode"`
	CoordinateX  float32   `json:"coordinate_x"`
	CoordinateY  float32   `json:"coordinate_y"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (*Address) GormDataType() string {
	return "jsonb"
}
