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
	Contacts       []entity.OfficeContact
	SocialNetworks []entity.OfficeSocialNetwork
}

type OfficeRepository interface {
	Create(context.Context, CreateOfficeData) (*entity.Office, error)
	UpdateByCouncillorID(context.Context, uint, UpdateOfficeData) (*entity.Office, error)
	FindByID(context.Context, uint) (*entity.Office, error)
	FindByCouncillorID(context.Context, uint) (*entity.Office, error)
}
