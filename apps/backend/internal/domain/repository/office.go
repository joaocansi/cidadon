package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateOfficeData struct {
	CouncillorID   uint
	Description    string
	Contacts       []entity.OfficeContact
	SocialNetworks []entity.OfficeSocialNetwork
}

type UpdateOfficeData struct {
	Description    string
	Contacts       []entity.OfficeContact
	SocialNetworks []entity.OfficeSocialNetwork
}

type OfficeRepository interface {
	Create(context.Context, CreateOfficeData) (*entity.Office, error)
	UpdateByCouncillorID(context.Context, uint, UpdateOfficeData) (*entity.Office, error)
	FindByID(context.Context, uint) (*entity.Office, error)
	FindByCouncillorID(context.Context, uint) (*entity.Office, error)
	ListByCityState(context.Context, string, string) ([]entity.Office, error)
	ListActiveOfficeIDsNear(context.Context, float64, float64, float64, int) ([]uint, error)
	Search(context.Context, string, string, string) ([]entity.Office, error)
}
