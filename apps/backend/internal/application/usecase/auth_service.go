package usecase

import (
	"cidadon/internal/adapters/external/provider"
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"cidadon/internal/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AuthServiceImpl struct {
	userRepo                repository.UserRepository
	sessionRepo             repository.SessionRepository
	citizenRepo             repository.CitizenRepository
	councillorRepo          repository.CouncillorRepository
	officeMemberRepo        repository.OfficeMemberRepository
	officeMemberRequestRepo repository.OfficeMemberRequestRepository
	transactionManager      database.TransactionManager
	jwtProvider             provider.JwtProvider
	hashProvider            provider.HashProvider
	cryptoProvider          provider.CryptoProvider
	addressEnricher         provider.AddressEnricher
	logger                  *zap.SugaredLogger
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	citizenRepo repository.CitizenRepository,
	councillorRepo repository.CouncillorRepository,
	officeMemberRepo repository.OfficeMemberRepository,
	officeMemberRequestRepo repository.OfficeMemberRequestRepository,
	transactionManager database.TransactionManager,
	jwtProvider provider.JwtProvider,
	hashProvider provider.HashProvider,
	cryptoProvider provider.CryptoProvider,
	addressEnricher provider.AddressEnricher,
	logger *zap.SugaredLogger) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:                userRepo,
		sessionRepo:             sessionRepo,
		citizenRepo:             citizenRepo,
		councillorRepo:          councillorRepo,
		officeMemberRepo:        officeMemberRepo,
		officeMemberRequestRepo: officeMemberRequestRepo,
		transactionManager:      transactionManager,
		jwtProvider:             jwtProvider,
		hashProvider:            hashProvider,
		cryptoProvider:          cryptoProvider,
		addressEnricher:         addressEnricher,
		logger:                  logger.Named("AuthService"),
	}
}

func (as *AuthServiceImpl) Login(ctx context.Context, data service.LoginInput) (*service.LoginOutput, error) {
	email := normalizeEmail(data.Email)
	user, err := as.userRepo.FindByEmail(ctx, email)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.Unauthorized("email/password does not match")
			}
			as.logger.Error("failed to find user", "email", email, "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	if ok := as.cryptoProvider.Compare(data.Password, user.Password); !ok {
		return nil, service.Unauthorized("email/password does not match")
	}

	subject := fmt.Sprintf("%d:%s", user.ID, user.Role)
	refreshToken, err := as.hashProvider.Generate()
	if err != nil {
		as.logger.Error("failed to generate refreshToken", "email", user.Email, "error", err)
		return nil, service.Internal(err)
	}

	accessToken, err := as.jwtProvider.Generate(subject)
	if err != nil {
		as.logger.Error("failed to generate accessToken", "email", user.Email, "error", err)
		return nil, service.Internal(err)
	}

	createSessionData := repository.CreateSessionData{
		UserID:           user.ID,
		RefreshTokenHash: refreshToken.Hash,
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 365),
		IpAddress:        "",
		UserAgent:        "",
	}

	createdSessionData, err := as.sessionRepo.Create(ctx, createSessionData)
	if err != nil {
		as.logger.Error("failed to create session", "email", user.Email, "error", err)
		return nil, service.Internal(err)
	}

	return &service.LoginOutput{
		RefreshTokenExpiresIn: createdSessionData.ExpiresAt,
		RefreshToken:          refreshToken.Value,
		AccessTokenExpiresIn:  accessToken.ExpiresAt,
		AccessToken:           accessToken.Token,
		Role:                  user.Role,
	}, nil
}

func (as *AuthServiceImpl) CurrentUser(ctx context.Context, userID uint) (*service.CurrentUserOutput, error) {
	user, err := as.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, service.NotFound("user not found")
	}
	output := &service.CurrentUserOutput{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role}
	if user.Role == entity.CouncillorUser {
		if councillor, findErr := as.councillorRepo.FindByUserID(ctx, userID); findErr == nil {
			output.ImageURL = councillor.ImageURL
		}
	}
	if user.Role == entity.OfficeMemberUser {
		if member, findErr := as.officeMemberRepo.FindByUserID(ctx, userID); findErr == nil {
			output.ImageURL = member.ImageURL
		}
	}
	return output, nil
}

func (as *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	if err := as.sessionRepo.RevokeByRefreshTokenHash(ctx, as.hashProvider.Hash(refreshToken)); err != nil {
		return service.Internal(err)
	}
	return nil
}

func (as *AuthServiceImpl) RegisterCitizen(ctx context.Context, input service.RegisterCitizenInput) (*service.RegisterCitizenOutput, error) {
	hashedPassword := as.cryptoProvider.Hash(input.Password)
	createCitizenData := repository.CreateCitizenData{
		Name:     strings.TrimSpace(input.Name),
		Email:    normalizeEmail(input.Email),
		Password: hashedPassword,
		City:     strings.TrimSpace(input.City),
		State:    strings.ToUpper(strings.TrimSpace(input.State)),
	}

	createdCitizen, err := as.citizenRepo.Create(ctx, createCitizenData)
	if err != nil {
		var dbError *repository.DBError
		if errors.As(err, &dbError) {
			if dbError.Code == repository.DBErrorConflict {
				return nil, service.Conflict("email already registered")
			}
			as.logger.Error("failed to create citizen", "email", input.Email, "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	return &service.RegisterCitizenOutput{
		RegisterBaseOutput: service.RegisterBaseOutput{
			Name:  createdCitizen.User.Name,
			Email: createdCitizen.User.Email,
		},
		City:  createdCitizen.City,
		State: createdCitizen.State,
	}, nil
}

func (as *AuthServiceImpl) RegisterCouncillor(ctx context.Context, input service.RegisterCouncillorInput) (*service.RegisterCouncillorOutput, error) {
	hashedPassword := as.cryptoProvider.Hash(input.Password)
	createCouncillorData := repository.CreateCouncillorData{
		Party:    strings.TrimSpace(input.Party),
		Name:     strings.TrimSpace(input.Name),
		Email:    normalizeEmail(input.Email),
		Password: hashedPassword,
		City:     strings.TrimSpace(input.City),
		State:    strings.ToUpper(strings.TrimSpace(input.State)),
		ImageURL: strings.TrimSpace(input.ImageURL),
	}

	createdCouncillor, err := as.councillorRepo.Create(ctx, createCouncillorData)
	if err != nil {
		var dbError *repository.DBError
		if errors.As(err, &dbError) {
			if dbError.Code == repository.DBErrorConflict {
				return nil, service.Conflict("email already registered")
			}
			as.logger.Error("failed to create councillor", "email", input.Email, "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	return &service.RegisterCouncillorOutput{
		RegisterBaseOutput: service.RegisterBaseOutput{
			Name:  createdCouncillor.User.Name,
			Email: createdCouncillor.User.Email,
		},
		ImageURL: createdCouncillor.ImageURL,
		Party:    createdCouncillor.Party,
		City:     createdCouncillor.City,
		State:    createdCouncillor.State,
	}, nil
}

func (as *AuthServiceImpl) RegisterOfficeMember(ctx context.Context, input service.RegisterOfficeMemberInput) (*service.RegisterOfficeMemberOutput, error) {
	officeMemberRequest, err := as.officeMemberRequestRepo.FindByToken(ctx, as.hashProvider.Hash(input.Token))
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.NotFound("token provided does not exist")
			}
			as.logger.Error("failed to find office member request", "token", utils.Mask(input.Token, 50), "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	hashedPassword := as.cryptoProvider.Hash(input.Password)
	createOfficeMemberData := repository.CreateOfficeMemberData{
		Email:    normalizeEmail(officeMemberRequest.Email),
		Password: hashedPassword,
		Name:     strings.TrimSpace(input.Name),
		OfficeID: officeMemberRequest.OfficeID,
		ImageURL: strings.TrimSpace(input.ImageURL),
	}

	var createdOfficeMember *entity.OfficeMember
	err = as.transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
		createdOfficeMember, err = as.officeMemberRepo.Create(ctx, createOfficeMemberData)
		if err != nil {
			var dbError *repository.DBError
			if errors.As(err, &dbError) && dbError.Code == repository.DBErrorConflict {
				return service.Conflict("email already registered")
			}
			as.logger.Error("failed to create office member", "email", officeMemberRequest.Email, "error", err)
			return service.Internal(err)
		}

		err = as.officeMemberRequestRepo.Delete(ctx, officeMemberRequest.ID)
		if err != nil {
			as.logger.Error("failed to delete office member request", "officeMemberRequestID", officeMemberRequest.ID, "error", err)
			return service.Internal(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &service.RegisterOfficeMemberOutput{
		RegisterBaseOutput: service.RegisterBaseOutput{
			Name:  createdOfficeMember.User.Name,
			Email: createdOfficeMember.User.Email,
		},
		OfficeID: createdOfficeMember.OfficeID,
		ImageURL: createdOfficeMember.ImageURL,
	}, nil
}

func (as *AuthServiceImpl) PreviewOfficeMemberInvitation(ctx context.Context, token string) (*service.OfficeMemberInvitationOutput, error) {
	request, err := as.officeMemberRequestRepo.FindByToken(ctx, as.hashProvider.Hash(strings.TrimSpace(token)))
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return nil, service.NotFound("invitation not found")
		}
		as.logger.Error("failed to preview office invitation", "error", err)
		return nil, service.Internal(err)
	}
	if request.Office == nil || request.Office.Councillor == nil || request.Office.Councillor.User == nil {
		as.logger.Error("office invitation is missing its public office data", "requestID", request.ID)
		return nil, service.Internal(fmt.Errorf("office invitation %d is missing office data", request.ID))
	}
	councillor := request.Office.Councillor
	return &service.OfficeMemberInvitationOutput{
		OfficeID:       request.OfficeID,
		CouncillorName: councillor.User.Name,
		Party:          councillor.Party,
		ImageURL:       councillor.ImageURL,
		City:           councillor.City,
		State:          councillor.State,
		ExpiresAt:      request.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
