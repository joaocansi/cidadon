package model

import (
	"cidadon/internal/domain/entity"
	"encoding/json"

	"gorm.io/gorm"
)

type Citizen struct {
	gorm.Model
	User    *User
	UserID  uint
	LivesIn json.RawMessage `database:"type:jsonb"`
}

func (c *Citizen) ToDomain() *entity.Citizen {
	var livesIn entity.Address
	_ = json.Unmarshal(c.LivesIn, &livesIn)

	return &entity.Citizen{
		ID:        c.ID,
		UserID:    c.UserID,
		User:      c.User.ToDomain(),
		LivesIn:   livesIn,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
