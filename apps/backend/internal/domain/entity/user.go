package entity

import (
	"time"
)

type UserRole string

const (
	CitizenUser      UserRole = "citizen"
	CouncillorUser   UserRole = "councillor"
	OfficeMemberUser UserRole = "office_member"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
