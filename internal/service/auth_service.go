package service

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/domain/service"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/provider"
	"cidadon/internal/utils"
	"context"
	"errors"
	"fmt"
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
	user, err := as.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.Unauthorized("email/password does not match")
			}
			as.logger.Error("failed to find user", "email", data.Email, "error", err)
			return nil, service.Internal(err)
		}
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
	}, nil
}

func (as *AuthServiceImpl) RegisterCitizen(ctx context.Context, input service.RegisterCitizenInput) (*service.RegisterCitizenOutput, error) {
	hashedPassword := as.cryptoProvider.Hash(input.Password)
	createCitizenData := repository.CreateCitizenData{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		City:     input.City,
		State:    input.State,
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
		Party:    input.Party,
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		City:     input.City,
		State:    input.State,
		ImageURL: input.ImageURL,
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
	officeMemberRequest, err := as.officeMemberRequestRepo.FindByToken(ctx, input.Token)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.NotFound("token provided does not exist")
			}
			as.logger.Error("failed to find office member request", "token", utils.Mask(input.Token, 50), "error", err)
			return nil, service.Internal(err)
		}
	}

	hashedPassword := as.cryptoProvider.Hash(input.Password)
	createOfficeMemberData := repository.CreateOfficeMemberData{
		Email:    officeMemberRequest.Email,
		Password: hashedPassword,
		Name:     input.Name,
		OfficeID: officeMemberRequest.OfficeID,
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
