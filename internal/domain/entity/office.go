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
	Contacts       []OfficeContact       `json:"contacts"`
	SocialNetworks []OfficeSocialNetwork `json:"social_networks"`
}

type OfficeContact struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Position int    `json:"position"`
}

type OfficeSocialNetwork struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Position int    `json:"position"`
}
