package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Role string

const (
	CitizenRole         Role = "citizen"
	CouncillorRole      Role = "councillor"
	CouncillorStaffRole Role = "councillor_staff"
)

func (u *User) ToModel() *Model {
	return &Model{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Role:     u.Role,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
