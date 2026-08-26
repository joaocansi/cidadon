package model

import (
	"cidadon/internal/domain/entity"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Role     string `model:"default:'citizen'"`
}

func (u *User) ToDomain() *entity.User {
	if u == nil {
		return nil
	}

	return &entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
