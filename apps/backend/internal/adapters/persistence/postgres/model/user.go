package model

import (
	"cidadon/internal/domain/entity"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:citizen;not null"`
}

func (u *User) ToDomain() *entity.User {
	if u == nil {
		return nil
	}
	return &entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      entity.UserRole(u.Role),
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
