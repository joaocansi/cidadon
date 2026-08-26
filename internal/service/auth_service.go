package service

import (
	"cidadon/internal/domain/repository"
	"cidadon/internal/domain/service"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/provider"
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, data LoginInput) (LoginOutput, error)
	RegisterCitizen(ctx context.Context, data RegisterCitizenInput) error
	RegisterCouncillor(ctx context.Context, data RegisterCouncillorInput) error
	RegisterOfficeMember(ctx context.Context, data RegisterOfficeMemberInput) error
}

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

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=12"`
}

type LoginOutput struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresIn  time.Time `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Time `json:"refresh_token_expires_in"`
}

func (as *AuthServiceImpl) Login(ctx context.Context, data LoginInput) (LoginOutput, error) {
	user, err := as.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		if stderrors.Is(err, repository.ErrDBNotFound) {
			return LoginOutput{}, service.NotFound("email/password does not match")
		}
		as.logger.Error("error finding user by email", err)
		return LoginOutput{}, service.Internal(err, "error finding user by email")
	}

	if ok := as.cryptoProvider.Compare(data.Password, user.Password); !ok {
		return LoginOutput{}, service.Unauthorized("email/password does not match")
	}

	subject := fmt.Sprintf("%d:%s", user.ID, user.Role)

	refreshToken, err := as.hashProvider.Generate()
	if err != nil {
		as.logger.Error("error generating refresh token", err)
		return LoginOutput{}, service.Internal(err, "error creating session")
	}

	accessToken, err := as.jwtProvider.Generate(subject)
	if err != nil {
		as.logger.Error("error generating access token", err)
		return LoginOutput{}, service.Internal(err, "error creating session")
	}

	createdSessionData := repository.CreateSessionData{
		UserID:           user.ID,
		RefreshTokenHash: refreshToken.Hash,
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 365),
		IpAddress:        "",
		UserAgent:        "",
	}

	_, err = as.sessionRepo.Create(ctx, createdSessionData)
	if err != nil {
		as.logger.Error("error creating session", err)
		return LoginOutput{}, service.Internal(err, "error creating session")
	}

	return LoginOutput{
		RefreshTokenExpiresIn: createdSessionData.ExpiresAt,
		RefreshToken:          refreshToken.Value,
		AccessTokenExpiresIn:  accessToken.ExpiresAt,
		AccessToken:           accessToken.Token,
	}, nil
}

type RegisterBaseInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=12"`
}

type RegisterCitizenInput struct {
	RegisterBaseInput
}

func (as *AuthServiceImpl) RegisterCitizen(ctx context.Context, input RegisterCitizenInput) error {
	user, err := as.userRepo.FindByEmail(ctx, input.Email)
	if user != nil {
		if err != nil {
			as.logger.Error("error finding user by email", err)
			return service.Internal(err, "error finding user by email")
		}
		return service.Conflict("email already registered")
	}

	hashedPassword := as.cryptoProvider.Hash(input.Password)
	if err != nil {
		as.logger.Error("error hashing password", err)
		return service.Internal(err, "error hashing password")
	}

	createCitizenData := repository.CreateCitizenData{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	_, err = as.citizenRepo.Create(ctx, createCitizenData)
	if err != nil {
		as.logger.Error("error creating citizen", err)
		return service.Internal(err, "error creating citizen")
	}

	return nil
}

type RegisterCouncillorInput struct {
	RegisterBaseInput
	Party    string `json:"party" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}

func (as *AuthServiceImpl) RegisterCouncillor(ctx context.Context, input RegisterCouncillorInput) error {
	user, err := as.userRepo.FindByEmail(ctx, input.Email)
	if user != nil {
		if err != nil {
			as.logger.Error("error finding user by email", err)
			return service.Internal(err, "error finding user by email")
		}
		return service.Conflict("email already registered")
	}

	hashedPassword := as.cryptoProvider.Hash(input.Password)
	if err != nil {
		as.logger.Error("error hashing password", err)
		return service.Internal(err, "error hashing password")
	}

	createCouncillorData := repository.CreateCouncillorData{
		Party:    input.Party,
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	_, err = as.councillorRepo.Create(ctx, createCouncillorData)
	if err != nil {
		as.logger.Error("error creating user", err)
		return service.Internal(err, "error creating user")
	}

	return nil
}

type RegisterOfficeMemberInput struct {
	RegisterBaseInput
	Token string `json:"token" binding:"required"`
}

func (as *AuthServiceImpl) RegisterOfficeMember(ctx context.Context, input RegisterOfficeMemberInput) error {
	_, err := as.officeMemberRequestRepo.FindByToken(ctx, input.Token)
	if err != nil {
		if stderrors.Is(err, repository.ErrDBNotFound) {
			return service.NotFound("no office member request found")
		}
		as.logger.Error("error finding office member", err)
		return service.Internal(err, "error finding office member")
	}

	user, err := as.userRepo.FindByEmail(ctx, input.Email)
	if user != nil {
		if err != nil {
			as.logger.Error("error finding user by email", err)
			return service.Internal(err, "error finding user by email")
		}
		return service.Conflict("email already registered")
	}

	hashedPassword := as.cryptoProvider.Hash(input.Password)
	if err != nil {
		as.logger.Error("error hashing password", err)
		return service.Internal(err, "error hashing password")
	}

	createOfficeMemberData := repository.CreateOfficeMemberData{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	}
	_, err = as.officeMemberRepo.Create(ctx, createOfficeMemberData)
	if err != nil {
		as.logger.Error("error creating user", err)
		return service.Internal(err, "error creating user")
	}

	return nil
}
