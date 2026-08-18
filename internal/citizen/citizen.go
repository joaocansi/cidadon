package citizen

import (
	"cidadon/internal/address"
	"cidadon/internal/user"
	"time"

	"gorm.io/gorm"
)

type Citizen struct {
	ID        uint
	User      *user.User
	UserID    uint
	LivesIn   address.Address
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Citizen) ToModel() *Model {
	if c == nil {
		return nil
	}

	return &Model{
		Model: gorm.Model{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
		UserID:  c.UserID,
		LivesIn: c.LivesIn,
		User:    c.User.ToModel(),
	}
}
