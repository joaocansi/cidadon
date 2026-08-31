package model

import (
	"cidadon/internal/domain/entity"
)

type Citizen struct {
	UserID uint   `gorm:"primaryKey"`
	User   *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	City   string `gorm:"not null"`
	State  string `gorm:"not null"`
}

func (c *Citizen) ToDomain() *entity.Citizen {
	return &entity.Citizen{
		UserID: c.UserID,
		User:   c.User.ToDomain(),
		City:   c.City,
		State:  c.State,
	}
}
