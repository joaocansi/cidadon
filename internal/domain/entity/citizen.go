package entity

import (
	"time"
)

type Citizen struct {
	ID        uint      `json:"id"`
	User      *User     `json:"user,omitempty"`
	UserID    uint      `json:"user_id"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
