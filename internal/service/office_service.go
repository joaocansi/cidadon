package service

import (
	"cidadon/internal/domain/repository"
	"cidadon/internal/domain/service"
	"cidadon/internal/provider"
	"context"
	"errors"

	"go.uber.org/zap"
)

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

func (s *OfficeServiceImpl) Create(ctx context.Context, input service.CreateOfficeInput) (*service.CreateOfficeOutput, error) {
	createOfficeData := repository.CreateOfficeData{
		CouncillorID:   input.CouncillorID,
		Contacts:       input.Contacts,
		SocialNetworks: input.SocialNetworks,
	}

	createdOffice, err := s.officeRepo.Create(ctx, createOfficeData)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorConflict {
				return nil, service.Conflict("office must be unique")
			}
			s.logger.Error("failed to create office member", "councillorID", input.CouncillorID, "error", err)
			return nil, service.Internal(err)
		}
	}

	return &service.CreateOfficeOutput{
		OfficeID:       createdOffice.ID,
		CouncillorID:   createdOffice.CouncillorID,
		Contacts:       createdOffice.Contacts,
		SocialNetworks: createdOffice.SocialNetworks,
	}, nil
}

func (s *OfficeServiceImpl) Update(ctx context.Context, input service.UpdateOfficeInput) (*service.UpdateOfficeOutput, error) {
	updateOfficeData := repository.UpdateOfficeData{
		Contacts:       input.Contacts,
		SocialNetworks: input.SocialNetworks,
	}
	updatedOffice, err := s.officeRepo.UpdateByCouncillorID(ctx, input.CouncillorID, updateOfficeData)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.NotFound("office not found")
			}
			s.logger.Error("failed to update office member", "councillorID", input.CouncillorID, "error", err)
			return nil, service.Internal(err)
		}
	}
	return &service.UpdateOfficeOutput{
		OfficeID:       updatedOffice.ID,
		CouncillorID:   updatedOffice.CouncillorID,
		Contacts:       updatedOffice.Contacts,
		SocialNetworks: updatedOffice.SocialNetworks,
	}, nil
}
