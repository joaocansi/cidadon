package usecase

import (
	"cidadon/internal/adapters/external/provider"
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

type OfficeServiceImpl struct {
	officeRepo              repository.OfficeRepository
	councillorRepo          repository.CouncillorRepository
	userRepo                repository.UserRepository
	officeMemberRepo        repository.OfficeMemberRepository
	officeMemberRequestRepo repository.OfficeMemberRequestRepository
	hashProvider            provider.HashProvider
	mailer                  provider.Mailer
	transactionManager      database.TransactionManager
	logger                  *zap.SugaredLogger
}

func (s *OfficeServiceImpl) InviteMember(ctx context.Context, userID uint, input service.InviteOfficeMemberInput) (*service.InviteOfficeMemberOutput, error) {
	office, err := s.officeRepo.FindByCouncillorID(ctx, userID)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return nil, service.NotFound("office not found")
		}
		return nil, service.Internal(err)
	}
	councillor, err := s.councillorRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, officeNotFoundOrInternal(err)
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if pending, findErr := s.officeMemberRequestRepo.FindByEmailAndOffice(ctx, email, office.ID); findErr == nil {
		if deleteErr := s.officeMemberRequestRepo.Delete(ctx, pending.ID); deleteErr != nil {
			return nil, service.Internal(deleteErr)
		}
	} else {
		var dbErr *repository.DBError
		if !errors.As(findErr, &dbErr) || dbErr.Code != repository.DBErrorNotFound {
			return nil, service.Internal(findErr)
		}
	}
	token, err := s.hashProvider.Generate()
	if err != nil {
		return nil, service.Internal(err)
	}
	expiresAt := time.Now().Add(72 * time.Hour)
	request, err := s.officeMemberRequestRepo.Create(ctx, repository.CreateOfficeMemberRequestData{OfficeID: office.ID, Email: email, Token: token.Hash, ExpiresAt: expiresAt})
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorConflict {
			return nil, service.Conflict("an invitation is already pending")
		}
		return nil, service.Internal(err)
	}
	inviteURL := fmt.Sprintf("%s/register/member?token=%s", frontendURL(), token.Value)
	invitation := provider.OfficeInvitation{
		CouncillorName: councillor.User.Name,
		Party:          councillor.Party,
		City:           councillor.City,
		State:          councillor.State,
	}
	subject := fmt.Sprintf("Convite para o Gabinete de %s no Cidadon", councillor.User.Name)
	if err := s.mailer.Send(email, subject, provider.OfficeInvitationHTML(inviteURL, invitation)); err != nil {
		if deleteErr := s.officeMemberRequestRepo.Delete(ctx, request.ID); deleteErr != nil {
			s.logger.Error("failed to remove undelivered invitation", "requestID", request.ID, "error", deleteErr)
		}
		s.logger.Error("failed to deliver office invitation", "email", email, "error", err)
		return nil, service.Unavailable("email delivery is unavailable")
	}
	return &service.InviteOfficeMemberOutput{Email: email, ExpiresAt: expiresAt.Format(time.RFC3339)}, nil
}

func (s *OfficeServiceImpl) ResolveOfficeID(ctx context.Context, userID uint, role entity.UserRole) (uint, error) {
	if role == entity.CouncillorUser {
		office, err := s.officeRepo.FindByCouncillorID(ctx, userID)
		if err != nil {
			return 0, service.NotFound("office not found")
		}
		return office.ID, nil
	}
	member, err := s.officeMemberRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, service.NotFound("office member not found")
	}
	return member.OfficeID, nil
}

func (s *OfficeServiceImpl) ListDirectory(ctx context.Context, city, state string) ([]service.OfficeDirectoryItem, error) {
	offices, err := s.officeRepo.ListByCityState(ctx, strings.TrimSpace(city), strings.ToUpper(strings.TrimSpace(state)))
	if err != nil {
		return nil, service.Internal(err)
	}
	items := make([]service.OfficeDirectoryItem, 0, len(offices))
	for _, office := range offices {
		item := service.OfficeDirectoryItem{OfficeID: office.ID}
		if office.Councillor != nil {
			item.Party = office.Councillor.Party
			if office.Councillor.User != nil {
				item.CouncillorName = office.Councillor.User.Name
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *OfficeServiceImpl) SearchPublic(ctx context.Context, query, city, state string) ([]service.PublicOfficeOutput, error) {
	offices, err := s.officeRepo.Search(ctx, strings.TrimSpace(query), strings.TrimSpace(city), strings.ToUpper(strings.TrimSpace(state)))
	if err != nil {
		return nil, service.Internal(err)
	}
	result := make([]service.PublicOfficeOutput, 0, len(offices))
	for _, office := range offices {
		result = append(result, publicOfficeOutput(office))
	}
	return result, nil
}

func (s *OfficeServiceImpl) FindPublic(ctx context.Context, id uint) (*service.PublicOfficeOutput, error) {
	office, err := s.officeRepo.FindByID(ctx, id)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return nil, service.NotFound("office not found")
		}
		return nil, service.Internal(err)
	}
	result := publicOfficeOutput(*office)
	return &result, nil
}

func (s *OfficeServiceImpl) FindManaged(ctx context.Context, userID uint, role entity.UserRole) (*service.ManagedOfficeOutput, error) {
	officeID, err := s.ResolveOfficeID(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	office, err := s.officeRepo.FindByID(ctx, officeID)
	if err != nil {
		return nil, service.NotFound("office not found")
	}
	members, err := s.officeMemberRepo.ListByOfficeID(ctx, officeID)
	if err != nil {
		return nil, service.Internal(err)
	}
	result := &service.ManagedOfficeOutput{PublicOfficeOutput: publicOfficeOutput(*office), Members: make([]service.OfficeMemberOutput, 0, len(members))}
	for _, member := range members {
		item := service.OfficeMemberOutput{UserID: member.UserID, ImageURL: member.ImageURL}
		if member.User != nil {
			item.Name = member.User.Name
			item.Email = member.User.Email
		}
		result.Members = append(result.Members, item)
	}
	return result, nil
}

func (s *OfficeServiceImpl) ListMemberRequests(ctx context.Context, councillorID uint) ([]service.OfficeMemberRequestOutput, error) {
	office, err := s.officeRepo.FindByCouncillorID(ctx, councillorID)
	if err != nil {
		return nil, officeNotFoundOrInternal(err)
	}
	requests, err := s.officeMemberRequestRepo.ListByOfficeID(ctx, office.ID)
	if err != nil {
		return nil, service.Internal(err)
	}
	result := make([]service.OfficeMemberRequestOutput, 0, len(requests))
	for _, request := range requests {
		result = append(result, service.OfficeMemberRequestOutput{
			ID:        request.ID,
			Email:     request.Email,
			ExpiresAt: request.ExpiresAt.Format(time.RFC3339),
			CreatedAt: request.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *OfficeServiceImpl) CancelMemberRequest(ctx context.Context, councillorID, requestID uint) error {
	office, err := s.officeRepo.FindByCouncillorID(ctx, councillorID)
	if err != nil {
		return officeNotFoundOrInternal(err)
	}
	request, err := s.officeMemberRequestRepo.FindByID(ctx, requestID)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return service.NotFound("invitation not found")
		}
		return service.Internal(err)
	}
	if request.OfficeID != office.ID {
		return service.Forbidden("invitation does not belong to this office")
	}
	if err := s.officeMemberRequestRepo.Delete(ctx, request.ID); err != nil {
		return service.Internal(err)
	}
	return nil
}

func (s *OfficeServiceImpl) RemoveMember(ctx context.Context, councillorID, memberID uint) error {
	office, err := s.officeRepo.FindByCouncillorID(ctx, councillorID)
	if err != nil {
		return officeNotFoundOrInternal(err)
	}
	member, err := s.officeMemberRepo.FindByUserID(ctx, memberID)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return service.NotFound("office member not found")
		}
		return service.Internal(err)
	}
	if member.OfficeID != office.ID {
		return service.Forbidden("member does not belong to this office")
	}
	if err := s.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if deleteErr := s.userRepo.DeleteByID(txCtx, member.UserID); deleteErr != nil {
			return service.Internal(deleteErr)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func publicOfficeOutput(office entity.Office) service.PublicOfficeOutput {
	result := service.PublicOfficeOutput{OfficeID: office.ID, Description: office.Description, Contacts: office.Contacts, SocialNetworks: office.SocialNetworks}
	if office.Councillor != nil {
		result.Party = office.Councillor.Party
		result.ImageURL = office.Councillor.ImageURL
		result.City = office.Councillor.City
		result.State = office.Councillor.State
		if office.Councillor.User != nil {
			result.CouncillorName = office.Councillor.User.Name
		}
	}
	return result
}

func NewOfficeService(
	officeRepo repository.OfficeRepository,
	councillorRepo repository.CouncillorRepository,
	userRepo repository.UserRepository,
	officeMemberRepo repository.OfficeMemberRepository,
	officeMemberRequestRepo repository.OfficeMemberRequestRepository,
	hashProvider provider.HashProvider,
	mailer provider.Mailer,
	transactionManager database.TransactionManager,
	logger *zap.SugaredLogger) *OfficeServiceImpl {
	return &OfficeServiceImpl{
		officeRepo:              officeRepo,
		councillorRepo:          councillorRepo,
		userRepo:                userRepo,
		officeMemberRepo:        officeMemberRepo,
		officeMemberRequestRepo: officeMemberRequestRepo,
		hashProvider:            hashProvider,
		mailer:                  mailer,
		transactionManager:      transactionManager,
		logger:                  logger.Named("OfficeService"),
	}
}

func frontendURL() string {
	if value := strings.TrimRight(strings.TrimSpace(os.Getenv("FRONTEND_URL")), "/"); value != "" {
		return value
	}
	return "http://localhost:3000"
}

func (s *OfficeServiceImpl) Create(ctx context.Context, input service.CreateOfficeInput) (*service.CreateOfficeOutput, error) {
	createOfficeData := repository.CreateOfficeData{
		CouncillorID:   input.CouncillorID,
		Description:    strings.TrimSpace(input.Description),
		Contacts:       normalizeContacts(input.Contacts),
		SocialNetworks: normalizeSocialNetworks(input.SocialNetworks),
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
		return nil, err
	}

	return &service.CreateOfficeOutput{
		OfficeID:       createdOffice.ID,
		CouncillorID:   createdOffice.CouncillorID,
		Description:    createdOffice.Description,
		Contacts:       createdOffice.Contacts,
		SocialNetworks: createdOffice.SocialNetworks,
	}, nil
}

func (s *OfficeServiceImpl) Update(ctx context.Context, input service.UpdateOfficeInput) (*service.UpdateOfficeOutput, error) {
	var updatedOffice *entity.Office
	var updatedCouncillor *entity.Councillor
	err := s.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var updateErr error
		updatedOffice, updateErr = s.officeRepo.UpdateByCouncillorID(txCtx, input.CouncillorID, repository.UpdateOfficeData{
			Description:    strings.TrimSpace(input.Description),
			Contacts:       normalizeContacts(input.Contacts),
			SocialNetworks: normalizeSocialNetworks(input.SocialNetworks),
		})
		if updateErr != nil {
			return officeNotFoundOrInternal(updateErr)
		}
		updatedCouncillor, updateErr = s.councillorRepo.UpdateByUserID(txCtx, input.CouncillorID, repository.UpdateCouncillorData{
			City:  strings.TrimSpace(input.City),
			State: strings.ToUpper(strings.TrimSpace(input.State)),
		})
		if updateErr != nil {
			return officeNotFoundOrInternal(updateErr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &service.UpdateOfficeOutput{
		OfficeID:       updatedOffice.ID,
		CouncillorID:   updatedOffice.CouncillorID,
		Description:    updatedOffice.Description,
		City:           updatedCouncillor.City,
		State:          updatedCouncillor.State,
		Contacts:       updatedOffice.Contacts,
		SocialNetworks: updatedOffice.SocialNetworks,
	}, nil
}

func officeNotFoundOrInternal(err error) error {
	var dbErr *repository.DBError
	if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
		return service.NotFound("office not found")
	}
	return service.Internal(err)
}

func normalizeContacts(items []entity.OfficeContact) []entity.OfficeContact {
	result := make([]entity.OfficeContact, 0, len(items))
	for _, item := range items {
		result = append(result, entity.OfficeContact{Type: strings.TrimSpace(item.Type), Value: strings.TrimSpace(item.Value), Position: len(result)})
	}
	return result
}

func normalizeSocialNetworks(items []entity.OfficeSocialNetwork) []entity.OfficeSocialNetwork {
	result := make([]entity.OfficeSocialNetwork, 0, len(items))
	for _, item := range items {
		result = append(result, entity.OfficeSocialNetwork{Type: strings.TrimSpace(item.Type), Value: strings.TrimSpace(item.Value), Position: len(result)})
	}
	return result
}
