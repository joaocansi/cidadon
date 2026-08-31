package entity

import (
	"time"
)

type Office struct {
	ID             uint                  `json:"id"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	CouncillorID   uint                  `json:"councillor_id"`
	Councillor     *Councillor           `json:"councillor,omitempty"`
	Description    string                `json:"description"`
	Contacts       []OfficeContact       `json:"contacts"`
	SocialNetworks []OfficeSocialNetwork `json:"social_networks"`
}

type OfficeContact struct {
	Type     string `json:"type" binding:"required,min=2,max=40"`
	Value    string `json:"value" binding:"required,min=2,max=2048"`
	Position int    `json:"position" binding:"gte=0"`
}

type OfficeSocialNetwork struct {
	Type     string `json:"type" binding:"required,min=2,max=40"`
	Value    string `json:"value" binding:"required,min=2,max=2048"`
	Position int    `json:"position" binding:"gte=0"`
}
