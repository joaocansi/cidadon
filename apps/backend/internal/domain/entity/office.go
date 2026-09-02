package entity

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type Office struct {
	ID             uint                  `json:"id"`
	Slug           string                `json:"slug"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	CouncillorID   uint                  `json:"councillor_id"`
	Councillor     *Councillor           `json:"councillor,omitempty"`
	Description    string                `json:"description"`
	Contacts       []OfficeContact       `json:"contacts"`
	SocialNetworks []OfficeSocialNetwork `json:"social_networks"`
}

// OfficeSlug produces a readable, stable public identifier. The councillor ID
// is part of the suffix so equal party/name combinations never conflict.
func OfficeSlug(party, councillorName string, councillorID uint) string {
	raw := norm.NFD.String(strings.ToLower(strings.TrimSpace(party + "-" + councillorName)))
	var result strings.Builder
	lastHyphen := false
	for _, character := range raw {
		if unicode.Is(unicode.Mn, character) {
			continue
		}
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			result.WriteRune(character)
			lastHyphen = false
			continue
		}
		if !lastHyphen && result.Len() > 0 {
			result.WriteByte('-')
			lastHyphen = true
		}
	}
	base := strings.Trim(result.String(), "-")
	if base == "" {
		base = "gabinete"
	}
	return fmt.Sprintf("%s-%d", base, councillorID)
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
