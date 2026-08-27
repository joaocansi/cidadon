package model

import (
	"cidadon/internal/domain/entity"
	"encoding/json"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	CouncillorID   uint            `gorm:"uniqueIndex;not null"`
	Councillor     *Councillor     `gorm:"foreignKey:CouncillorID;references:UserID"`
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
		CouncillorID:   o.CouncillorID,
		Councillor:     o.Councillor.ToDomain(),
		Contacts:       Contacts,
		SocialNetworks: socialNetworks,
	}
}
