package bootstrap

import (
	"cidadon/internal/adapters/external/provider"
	"cidadon/internal/adapters/http/handler"
	http "cidadon/internal/adapters/http/middleware"
	"cidadon/internal/adapters/http/routes"
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/adapters/persistence/postgres/repository"
	"cidadon/internal/application/usecase"
	"cidadon/internal/domain/entity"
	environment "cidadon/internal/platform/config"
	"cidadon/internal/platform/database"
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
		&model.OfficeMemberRequest{},
		&model.Demand{},
		&model.DemandAssignment{},
		&model.DemandEvent{},
		&model.DemandComment{},
		&model.DemandCommentReport{},
		&model.Notification{},
	)
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
	if err := db.Model(&model.Demand{}).Where("status = ?", "seen").Update("status", string(entity.DemandStatusReview)).Error; err != nil {
		return err
	}
	if err := db.Model(&model.Demand{}).Where("status = ?", "resolved").Update("status", string(entity.DemandStatusCompleted)).Error; err != nil {
		return err
	}
	var legacyDemands []model.Demand
	if err := db.Where("NOT EXISTS (SELECT 1 FROM demand_events WHERE demand_events.demand_id = demands.id)").Find(&legacyDemands).Error; err != nil {
		return err
	}
	for _, demand := range legacyDemands {
		if err := db.Create(&model.DemandEvent{DemandID: demand.ID, Type: "migrated", Metadata: []byte(`{"source":"legacy_status"}`)}).Error; err != nil {
			return err
		}
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
	demandRepository := repository.NewDemandRepository(baseRepository)
	transactionManager := database.NewTransactionManager(db)

	authService := usecase.NewAuthService(userRepository, sessionRepository, citizenRepository, councillorRepository, officeMemberRepository, officeMemberRequestRepository, transactionManager, jwtProvider, hashProvider, cryptoProvider, nominatimAddressEnricher, logger)
	authHandler := handler.NewAuthHandler(authService)

	officeService := usecase.NewOfficeService(officeRepository, councillorRepository, userRepository, officeMemberRepository, officeMemberRequestRepository, hashProvider, provider.NewSMTPMailer(logger), transactionManager, logger)
	officeHandler := handler.NewOfficeHandler(officeService)
	demandService := usecase.NewDemandService(demandRepository, officeRepository, db, logger)
	demandHandler := handler.NewDemandHandler(demandService, officeService)

	router := gin.New()
	router.Use(http.ErrorHandler())

	routes.RegisterAuth(router, authMiddleware, authHandler)
	routes.RegisterOffice(router, authMiddleware, officeHandler)
	routes.RegisterDemands(router, authMiddleware, demandHandler)
	routes.RegisterNotifications(router, authMiddleware, demandHandler)

	if err := netHttp.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
