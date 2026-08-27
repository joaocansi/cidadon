package internal

import (
	"cidadon/internal/api/http"
	"cidadon/internal/domain/entity"
	"cidadon/internal/environment"
	"cidadon/internal/handler"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"cidadon/internal/provider"
	"cidadon/internal/repository"
	"cidadon/internal/service"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Citizen{},
		&model.Councillor{},
		&model.Office{},
		&model.OfficeMember{},
		&model.OfficeMemberRequest{})
}

func Run() error {
	if err := environment.Load(); err != nil {
		return err
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	rsa256Pair, err := environment.LoadRsaKeys(environment.Env.JwtProvider.PrivateKey, environment.Env.JwtProvider.PublicKey)
	if err != nil {
		return err
	}

	db, err := database.NewPostgresConnection(environment.Env.Database)
	if err != nil {
		return err
	}

	if err := migrate(db); err != nil {
		return err
	}

	hashProvider := provider.NewHashProvider()
	cryptoProvider := provider.NewCryptoProvider()
	jwtProvider := provider.NewJwtProvider(
		rsa256Pair.PublicKey,
		rsa256Pair.PrivateKey,
		environment.Env.JwtProvider.Issuer,
		environment.Env.JwtProvider.Audience,
		environment.Env.JwtProvider.Expiration)

	authMiddleware := http.NewAuthMiddleware(jwtProvider)

	nominatimAddressEnricher := provider.NewNominatimAddressEnricher(&netHttp.Client{}, "https://nominatim.openstreetmap.org/")

	baseRepository := database.NewBaseRepository(db)
	councillorRepository := repository.NewCouncillorRepository(baseRepository)
	citizenRepository := repository.NewCitizenRepository(baseRepository)
	userRepository := repository.NewUserRepository(baseRepository)
	sessionRepository := repository.NewSessionRepository(baseRepository)
	officeRepository := repository.NewOfficeRepository(baseRepository)
	officeMemberRepository := repository.NewOfficeMemberRepository(baseRepository)
	officeMemberRequestRepository := repository.NewOfficeMemberRequestRepository(baseRepository)
	transactionManager := database.NewTransactionManager(db)

	authService := service.NewAuthService(userRepository, sessionRepository, citizenRepository, councillorRepository, officeMemberRepository, officeMemberRequestRepository, transactionManager, jwtProvider, hashProvider, cryptoProvider, nominatimAddressEnricher, logger)
	authHandler := handler.NewAuthHandler(authService)

	officeService := service.NewOfficeService(officeRepository, officeMemberRepository, officeMemberRequestRepository, hashProvider, logger)
	officeHandler := handler.NewOfficeHandler(officeService)

	router := gin.New()
	router.Use(http.ErrorHandler())

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", authHandler.Login)
		authRouter.POST("/register/citizen", authHandler.RegisterCitizen)
		authRouter.POST("/register/councillor", authHandler.RegisterCouncillor)
		authRouter.POST("/register/office-member", authHandler.RegisterOfficeMember)
	}

	officeRouter := router.Group("/office")
	{
		onlyCouncillor := officeRouter.Group("/")
		onlyCouncillor.Use(authMiddleware.AuthHandler(entity.CouncillorUser))
		{
			onlyCouncillor.POST("/", officeHandler.Create)
			onlyCouncillor.POST("/member-request", officeHandler.NewMemberRequest)
		}

		councillorOrMember := officeRouter.Group("/")
		councillorOrMember.Use(authMiddleware.AuthHandler(entity.CouncillorUser, entity.OfficeMemberUser))
		{
			councillorOrMember.PUT("/", officeHandler.Update)
		}
	}

	if err := netHttp.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
