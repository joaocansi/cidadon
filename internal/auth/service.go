package auth

import (
	"cidadon/internal/address"
	"cidadon/internal/app/database"
	apperrors "cidadon/internal/app/errors"
	"cidadon/internal/app/providers"
	"cidadon/internal/auth/session"
	"cidadon/internal/citizen"
	"cidadon/internal/councillor"
	"cidadon/internal/user"
	"context"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, data LoginInput) (LoginOutput, error)
	Register(ctx context.Context, data RegisterInput) error
}

type ServiceImpl struct {
	userRepo             user.Repository
	sessionRepo          session.Repository
	citizenRepo          citizen.Repository
	councillorRepo       councillor.Repository
	transactionManager   database.TransactionManager
	jwtProvider          providers.JwtProvider
	refreshTokenProvider providers.RefreshTokenProvider
	cryptoProvider       providers.CryptoProvider
	addressEnricher      address.Enricher
	logger               *zap.SugaredLogger
}

func NewAuthService(
	userRepo user.Repository,
	sessionRepo session.Repository,
	citizenRepo citizen.Repository,
	councillorRepo councillor.Repository,
	transactionManager database.TransactionManager,
	jwtProvider providers.JwtProvider,
	refreshTokenProvider providers.RefreshTokenProvider,
	cryptoProvider providers.CryptoProvider,
	addressEnricher address.Enricher,
	logger *zap.SugaredLogger) *ServiceImpl {
	return &ServiceImpl{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		citizenRepo:          citizenRepo,
		councillorRepo:       councillorRepo,
		transactionManager:   transactionManager,
		jwtProvider:          jwtProvider,
		refreshTokenProvider: refreshTokenProvider,
		cryptoProvider:       cryptoProvider,
		addressEnricher:      addressEnricher,
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
		if stderrors.Is(err, apperrors.ErrDBNotFound) {
			return LoginOutput{}, apperrors.NotFound("email/password does not match")
		}
		as.logger.Error("error finding user by email", err)
		return LoginOutput{}, err
	}

	if ok := as.cryptoProvider.Compare(data.Password, findUser.Password); !ok {
		return LoginOutput{}, apperrors.Unauthorized("email/password does not match")
	}

	subject := fmt.Sprintf("%d:%s", findUser.ID, findUser.Role)
	refreshToken, err := as.refreshTokenProvider.Generate()
	accessToken, err := as.jwtProvider.Generate(subject)

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

type SubjectAuth struct {
	UserID uint      `json:"userId"`
	Role   user.Role `json:"role"`
}

func GetSubjectAuth(subject string) (*SubjectAuth, error) {
	subjectSplit := strings.Split(subject, ":")
	if len(subjectSplit) != 2 {
		return nil, fmt.Errorf("invalid subject: %s", subject)
	}

	userID, err := strconv.ParseUint(subjectSplit[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid subject: %s", subject)
	}

	return &SubjectAuth{
		Role:   user.Role(subjectSplit[1]),
		UserID: uint(userID),
	}, nil
}

type RegionCitizenInput struct {
	CoordinateX float32 `json:"coordinate_x" binding:"required"`
	CoordinateY float32 `json:"coordinate_y" binding:"required"`
}

type RegisterCitizenInput struct {
	Region RegionCitizenInput `json:"address" binding:"required"`
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
			return apperrors.InvalidInput("missing citizen data for citizen role")
		}
	case user.CouncillorRole:
		if data.Councillor == nil {
			return apperrors.InvalidInput("missing councillor data for councillor role")
		}
	default:
		return apperrors.InvalidInput("invalid role: %s", data.Role)
	}

	findUser, _ := as.userRepo.FindByEmail(ctx, data.Email)
	if findUser != nil {
		return apperrors.AlreadyExists("email already exists")
	}

	hashedPassword := as.cryptoProvider.Hash(data.Password)
	newUser := &user.User{
		Password: hashedPassword,
		Role:     data.Role,
		Email:    data.Email,
		Name:     data.Name,
	}

	registerLogger := as.logger.With("Register", newUser.Email)
	if data.Citizen != nil {
		userAddress := &address.Address{
			CoordinateX: data.Citizen.Region.CoordinateX,
			CoordinateY: data.Citizen.Region.CoordinateY,
		}

		err := as.addressEnricher.Enrich(ctx, userAddress)
		if err != nil {
			registerLogger.Error("error enriching address", err)
			return err
		}

		err = as.citizenRepo.Create(ctx, &citizen.Citizen{
			LivesIn: *userAddress,
			User:    newUser,
		})

		if err != nil {
			registerLogger.Error("error creating citizen", err)
			return err
		}

		return nil
	}

	if data.Councillor != nil {
		err := as.councillorRepo.Create(ctx, &councillor.Councillor{
			Party: data.Councillor.Party,
		})

		if err != nil {
			registerLogger.Error("error creating councillor", err)
			return err
		}

		return nil
	}

	return apperrors.InvalidInput("no councillor or citizen in RegisterInput")
}
