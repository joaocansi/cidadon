package user

import (
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Role     Role `model:"default:'citizen'"`
}

func (u *Model) ToDomain() *User {
	return &User{
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
