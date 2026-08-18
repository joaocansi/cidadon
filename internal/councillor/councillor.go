package councillor

import (
	"time"

	"gorm.io/gorm"
)

type Councillor struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Party     string
}

func (councillor *Councillor) ToModel() *Model {
	return &Model{
		Model: gorm.Model{
			ID:        councillor.ID,
			CreatedAt: councillor.CreatedAt,
			UpdatedAt: councillor.UpdatedAt,
		},
		Party: councillor.Party,
	}
}
