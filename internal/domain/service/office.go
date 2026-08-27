package service

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateOfficeInput struct {
	CouncillorID   uint
	Contacts       []entity.OfficeContact       `json:"contacts" binding:"required"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks" binding:"required"`
}

type CreateOfficeOutput struct {
	OfficeID       uint                         `json:"office_id"`
	CouncillorID   uint                         `json:"councillor_id"`
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}

type UpdateOfficeInput struct {
	CouncillorID   uint
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}

type UpdateOfficeOutput struct {
	OfficeID       uint                         `json:"office_id"`
	CouncillorID   uint                         `json:"councillor_id"`
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}

type OfficeService interface {
	Create(context.Context, CreateOfficeInput) (*CreateOfficeOutput, error)
	Update(context.Context, UpdateOfficeInput) (*UpdateOfficeOutput, error)
}
