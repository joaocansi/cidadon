package domain

import (
	"time"
)

type OfficeMember struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
	UserID    uint      `json:"user_id"`
	OfficeID  uint      `json:"office_id"`
	Office    *Office   `json:"office,omitempty"`
	ImageURL  string    `json:"image_url"`
}
