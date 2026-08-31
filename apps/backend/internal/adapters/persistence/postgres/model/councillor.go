package model

import (
	"cidadon/internal/domain/entity"
)

type Councillor struct {
	UserID   uint   `gorm:"primaryKey"`
	User     *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Party    string `gorm:"not null"`
	ImageURL string
	State    string `gorm:"not null"`
	City     string `gorm:"not null"`
}

func (c *Councillor) ToDomain() *entity.Councillor {
	if c == nil {
		return nil
	}

	return &entity.Councillor{
		UserID:   c.UserID,
		User:     c.User.ToDomain(),
		Party:    c.Party,
		ImageURL: c.ImageURL,
		City:     c.City,
		State:    c.State,
	}
}
