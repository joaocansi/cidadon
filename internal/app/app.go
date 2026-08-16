package app

import (
	"cidadon/internal/app/shared/database"
	"cidadon/internal/app/shared/security"
	"cidadon/internal/app/shared/settings"
	"cidadon/internal/app/shared/utils"
	"cidadon/internal/auth"
	"cidadon/internal/auth/session"
	"cidadon/internal/citizen"
	"cidadon/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&user.Model{}, &session.Model{}, &citizen.Model{})
}

func Run() error {
	if err := settings.Load(); err != nil {
		return err
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	rsa256Pair, err := utils.LoadRSAKeys(settings.Env.JwtProvider.PrivateKey, settings.Env.JwtProvider.PublicKey)
	if err != nil {
		return err
	}

	db, err := database.NewPostgresConnection(settings.Env.Database)
	if err != nil {
		return err
	}

	if err := migrate(db); err != nil {
		return err
	}

	refreshTokenProvider := security.NewRefreshTokenProvider(settings.Env.RefreshTokenProvider.Expiration)
	cryptoProvider := security.NewCryptoProvider()
	jwtProvider := security.NewJwtProvider(
		rsa256Pair.PublicKey,
		rsa256Pair.PrivateKey,
		settings.Env.JwtProvider.Issuer,
		settings.Env.JwtProvider.Audience,
		settings.Env.JwtProvider.Expiration)

	baseRepository := database.NewBaseRepository(db)
	userRepository := user.NewUserRepository(baseRepository, logger)
	sessionRepository := session.NewSessionRepository(baseRepository, logger)
	transactionManager := database.NewTransactionManager(db)

	authService := auth.NewAuthService(userRepository, sessionRepository, transactionManager, jwtProvider, refreshTokenProvider, cryptoProvider, logger)
	authHandler := auth.NewAuthHandler(authService)

	router := gin.New()
	router.Use(ErrorHandler())

	authRouter := router.Group("/auth")
	{
		authRouter.GET("/login", authHandler.Login)
		authRouter.POST("/register", authHandler.Register)
	}

	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
