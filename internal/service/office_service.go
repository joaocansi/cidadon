package service

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/domain/service"
	"cidadon/internal/provider"
	"context"
	"errors"

	"go.uber.org/zap"
)

type CreateOfficeInput struct {
	CouncillorID   uint
	Contacts       []entity.OfficeContact       `json:"contacts" binding:"required"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks" binding:"required"`
}

type CreateOfficeOutput struct{}

type UpdateOfficeInput struct {
	CouncillorID   uint
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}

type UpdateOfficeOutput struct{}

type OfficeService interface {
	Create(context.Context, CreateOfficeInput) (*CreateOfficeOutput, error)
}

type OfficeServiceImpl struct {
	officeRepo              repository.OfficeRepository
	officeMemberRepo        repository.OfficeMemberRepository
	officeMemberRequestRepo repository.OfficeMemberRequestRepository
	hashProvider            provider.HashProvider
	logger                  *zap.SugaredLogger
}

func NewOfficeService(
	officeRepo repository.OfficeRepository,
	officeMemberRepo repository.OfficeMemberRepository,
	officeMemberRequestRepo repository.OfficeMemberRequestRepository,
	hashProvider provider.HashProvider,
	logger *zap.SugaredLogger) *OfficeServiceImpl {
	return &OfficeServiceImpl{
		officeRepo:              officeRepo,
		officeMemberRepo:        officeMemberRepo,
		officeMemberRequestRepo: officeMemberRequestRepo,
		hashProvider:            hashProvider,
		logger:                  logger.Named("OfficeService"),
	}
}

func (s *OfficeServiceImpl) Create(ctx context.Context, input CreateOfficeInput) (*CreateOfficeOutput, error) {
	findOffice, err := s.officeRepo.FindByCouncillorID(ctx, input.CouncillorID)
	if findOffice != nil {
		return nil, service.Conflict("councillor already created his office")
	}

	if err != nil {
		return nil, service.Internal(err, "failed to find office")
	}

	createOfficeData := repository.CreateOfficeData{
		CouncillorID:   input.CouncillorID,
		Contacts:       input.Contacts,
		SocialNetworks: input.SocialNetworks,
	}

	_, err = s.officeRepo.Create(ctx, createOfficeData)
	if err != nil {
		return nil, service.Internal(err, "failed to create office")
	}

	return &CreateOfficeOutput{}, nil
}

func (s *OfficeServiceImpl) Update(ctx context.Context, input UpdateOfficeInput) (*UpdateOfficeOutput, error) {
	findOffice, err := s.officeRepo.FindByCouncillorID(ctx, input.CouncillorID)
	if err != nil {
		if errors.Is(err, repository.ErrDBNotFound) {
			return nil, repository.WrapDB(repository.ErrDBNotFound, "office not found")
		}
		return nil, service.NotFound("councillor already created his office")
	}
	updateOfficeData := repository.UpdateOfficeData{
		OfficeID:       findOffice.ID,
		Contacts:       input.Contacts,
		SocialNetworks: input.SocialNetworks,
	}
	err = s.officeRepo.Update(ctx, updateOfficeData)
	if err != nil {
		return nil, service.Internal(err, "failed to update office")
	}
	return nil, nil
}
