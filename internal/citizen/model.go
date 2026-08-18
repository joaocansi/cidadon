package citizen

import (
	"cidadon/internal/address"
	"cidadon/internal/user"

	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	User    *user.Model
	UserID  uint
	LivesIn address.Address
}

func (*Model) TableName() string {
	return "citizens"
}

func (m *Model) ToDomain() *Citizen {
	return &Citizen{
		ID:        m.ID,
		UpdatedAt: m.UpdatedAt,
		CreatedAt: m.CreatedAt,
		LivesIn:   m.LivesIn,
	}
}
