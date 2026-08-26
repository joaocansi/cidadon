package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateOfficeData struct {
	CouncillorID   uint
	Contacts       []entity.OfficeContact
	SocialNetworks []entity.OfficeSocialNetwork
}

type UpdateOfficeData struct {
	OfficeID       uint
	Contacts       []entity.OfficeContact
	SocialNetworks []entity.OfficeSocialNetwork
}

type OfficeRepository interface {
	Create(context.Context, CreateOfficeData) (*entity.Office, error)
	Update(context.Context, UpdateOfficeData) error
	FindByID(context.Context, uint) (*entity.Office, error)
	FindByCouncillorID(context.Context, uint) (*entity.Office, error)
}
