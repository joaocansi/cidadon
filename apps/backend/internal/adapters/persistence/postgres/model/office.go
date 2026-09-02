package model

import (
	"cidadon/internal/domain/entity"
	"encoding/json"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	CouncillorID   uint            `gorm:"uniqueIndex;not null"`
	Slug           string          `gorm:"uniqueIndex"`
	Councillor     *Councillor     `gorm:"foreignKey:CouncillorID;references:UserID"`
	Description    string          `gorm:"type:text"`
	Contacts       json.RawMessage `gorm:"type:jsonb"`
	SocialNetworks json.RawMessage `gorm:"type:jsonb"`
}

func (o *Office) ToDomain() *entity.Office {
	if o == nil {
		return nil
	}

	socialNetworks := make([]entity.OfficeSocialNetwork, 0)
	_ = json.Unmarshal(o.SocialNetworks, &socialNetworks)

	Contacts := make([]entity.OfficeContact, 0)
	_ = json.Unmarshal(o.Contacts, &Contacts)

	return &entity.Office{
		ID:             o.ID,
		Slug:           o.Slug,
		CouncillorID:   o.CouncillorID,
		Councillor:     o.Councillor.ToDomain(),
		Description:    o.Description,
		Contacts:       Contacts,
		SocialNetworks: socialNetworks,
	}
}
