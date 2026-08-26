package domain

import (
	"time"
)

type Councillor struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
	UserID    uint      `json:"user_id"`
	Party     string    `json:"party"`
	ImageURL  string    `json:"image_url"`
}
