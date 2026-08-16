package region

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	City         string
	State        string
	Neighborhood string
	Postcode     string
	CoordinateX  float32
	CoordinateY  float32
}

func (m Model) TableName() string {
	return "regions"
}
