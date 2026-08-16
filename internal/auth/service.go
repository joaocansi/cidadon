package auth

import (
	"cidadon/internal/app/shared/database"
	"cidadon/internal/app/shared/security"
	"cidadon/internal/auth/session"
	"cidadon/internal/user"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, data LoginInput) (LoginOutput, error)
	Register(ctx context.Context, data RegisterInput) error
}

type ServiceImpl struct {
	userRepo             user.Repository
	sessionRepo          session.Repository
	transactionManager   database.TransactionManager
	jwtProvider          security.JwtProvider
	refreshTokenProvider security.RefreshTokenProvider
	cryptoProvider       security.CryptoProvider
	logger               *zap.SugaredLogger
}

func NewAuthService(
	userRepo user.Repository,
	sessionRepo session.Repository,
	transactionManager database.TransactionManager,
	jwtProvider security.JwtProvider,
	refreshTokenProvider security.RefreshTokenProvider,
	cryptoProvider security.CryptoProvider,
	logger *zap.SugaredLogger) *ServiceImpl {
	return &ServiceImpl{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		transactionManager:   transactionManager,
		jwtProvider:          jwtProvider,
		refreshTokenProvider: refreshTokenProvider,
		cryptoProvider:       cryptoProvider,
		logger:               logger.Named("AuthService"),
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

func (as *ServiceImpl) Login(ctx context.Context, data LoginInput) (LoginOutput, error) {
	findUser, err := as.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return LoginOutput{}, nil
		}
		logger.Error("error finding user by email", err)
		return LoginOutput{}, err
	}

	if ok := as.cryptoProvider.Compare(data.Password, findUser.Password); !ok {
		return LoginOutput{}, errors.New(fmt.Sprintf("email/password does not match"))
	}

	userID := strconv.FormatUint(uint64(findUser.ID), 10)
	refreshToken, err := as.refreshTokenProvider.Generate()
	accessToken, err := as.jwtProvider.Generate(userID)

	createdSession := session.Session{
		UserID:           findUser.ID,
		RefreshTokenHash: refreshToken.Hash,
		ExpiresAt:        refreshToken.ExpiresAt,
	}

	err = as.sessionRepo.Create(ctx, createdSession)
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		RefreshTokenExpiresIn: refreshToken.ExpiresAt,
		RefreshToken:          refreshToken.Value,
		AccessTokenExpiresIn:  accessToken.ExpiresAt,
		AccessToken:           accessToken.Token,
	}, nil
}

type RegionCitizenInput struct {
	CoordinateX float32 `json:"coordinate_x" binding:"required"`
	CoordinateY float32 `json:"coordinate_y" binding:"required"`
}

type RegisterCitizenInput struct {
	Region RegionCitizenInput `json:"region" binding:"required"`
}

type RegisterCouncillorInput struct {
	Party string `json:"party" binding:"required"`
}

type RegisterInput struct {
	Name       string                   `json:"name" binding:"required"`
	Email      string                   `json:"email" binding:"required,email"`
	Password   string                   `json:"password" binding:"required,gte=6,lte=12"`
	Role       user.Role                `json:"role" binding:"required,oneof=citizen councillor"`
	Citizen    *RegisterCitizenInput    `json:"citizen"`
	Councillor *RegisterCouncillorInput `json:"councillor"`
}

func (as *ServiceImpl) Register(ctx context.Context, data RegisterInput) error {
	switch data.Role {
	case user.CitizenRole:
		if data.Citizen == nil {
			return errors.New("no citizen in RegisterInput")
		}
	case user.CouncillorRole:
		if data.Councillor == nil {
			return errors.New("no councillor in RegisterInput")
		}
	default:
		return errors.New("invalid role in RegisterInput")
	}
	return nil
}
