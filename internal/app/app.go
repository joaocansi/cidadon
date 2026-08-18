package app

import (
	"cidadon/internal/address"
	"cidadon/internal/app/database"
	"cidadon/internal/app/environment"
	"cidadon/internal/app/http"
	security2 "cidadon/internal/app/providers"
	"cidadon/internal/app/utils"
	"cidadon/internal/auth"
	"cidadon/internal/auth/session"
	"cidadon/internal/citizen"
	"cidadon/internal/councillor"
	"cidadon/internal/user"
	netHttp "net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&user.Model{}, &session.Model{}, &citizen.Model{}, &councillor.Model{})
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

	refreshTokenProvider := security2.NewRefreshTokenProvider(environment.Env.RefreshTokenProvider.Expiration)
	cryptoProvider := security2.NewCryptoProvider()
	jwtProvider := security2.NewJwtProvider(
		rsa256Pair.PublicKey,
		rsa256Pair.PrivateKey,
		environment.Env.JwtProvider.Issuer,
		environment.Env.JwtProvider.Audience,
		environment.Env.JwtProvider.Expiration)

	authMiddleware := http.NewAuthMiddleware(jwtProvider)

	nominatimAddressEnricher := address.NewNominatimAddressEnricher(&netHttp.Client{}, "https://nominatim.openstreetmap.org/")

	baseRepository := database.NewBaseRepository(db)
	councillorRepository := councillor.NewCouncillorRepository(baseRepository, logger)
	citizenRepository := citizen.NewCitizenRepository(baseRepository, logger)
	userRepository := user.NewUserRepository(baseRepository, logger)
	sessionRepository := session.NewSessionRepository(baseRepository, logger)
	transactionManager := database.NewTransactionManager(db)

	authService := auth.NewAuthService(userRepository, sessionRepository, citizenRepository, councillorRepository, transactionManager, jwtProvider, refreshTokenProvider, cryptoProvider, nominatimAddressEnricher, logger)
	authHandler := auth.NewAuthHandler(authService)

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
