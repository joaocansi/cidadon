package api

import (
	"cidadon/internal/app/database"
	"cidadon/internal/app/http"
	"cidadon/internal/app/utils"
	"cidadon/internal/environment"
	"cidadon/internal/handler"
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
		&model.OfficeContact{},
		&model.OfficeSocialNetwork{})
}

func Run() error {
	if err := environment.Load(); err != nil {
		return err
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	rsa256Pair, err := utils.LoadRSAKeys(environment.Env.JwtProvider.PrivateKey, environment.Env.JwtProvider.PublicKey)
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

	refreshTokenProvider := provider.NewRefreshTokenProvider(environment.Env.RefreshTokenProvider.Expiration)
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
	councillorRepository := repository.NewCouncillorRepository(baseRepository, logger)
	citizenRepository := repository.NewCitizenRepository(baseRepository, logger)
	userRepository := repository.NewUserRepository(baseRepository, logger)
	sessionRepository := repository.NewSessionRepository(baseRepository, logger)
	transactionManager := database.NewTransactionManager(db)

	authService := service.NewAuthService(userRepository, sessionRepository, citizenRepository, councillorRepository, transactionManager, jwtProvider, refreshTokenProvider, cryptoProvider, nominatimAddressEnricher, logger)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.New()
	router.Use(http.ErrorHandler())

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", authHandler.Login)
		authRouter.POST("/register", authHandler.Register)
		authRouter.GET("/profile", authMiddleware.AuthHandler(), authHandler.Profile)
	}

	if err := netHttp.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
