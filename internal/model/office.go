package model

import (
	"cidadon/internal/domain/entity"
	"encoding/json"

	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	CouncillorID   uint
	Councillor     *Councillor
	Contacts       json.RawMessage `database:"type:jsonb"`
	SocialNetworks json.RawMessage `database:"type:jsonb"`
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
		CouncillorID:   o.CouncillorID,
		Councillor:     o.Councillor.ToDomain(),
		Contacts:       Contacts,
		SocialNetworks: socialNetworks,
	}
}
