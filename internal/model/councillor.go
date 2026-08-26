package model

import (
	"cidadon/internal/domain/entity"

	"gorm.io/gorm"
)

type Councillor struct {
	gorm.Model
	UserID   uint
	User     *User
	Party    string
	ImageURL string
	State    string
	City     string
}

func (c *Councillor) ToDomain() *entity.Councillor {
	if c == nil {
		return nil
	}

	return &entity.Councillor{
		ID:        c.ID,
		UserID:    c.UserID,
		User:      c.User.ToDomain(),
		Party:     c.Party,
		ImageURL:  c.ImageURL,
		City:      c.City,
		State:     c.State,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
