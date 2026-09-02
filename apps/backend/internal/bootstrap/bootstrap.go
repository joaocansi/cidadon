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
	platformmedia "cidadon/internal/platform/media"
	"fmt"
	netHttp "net/http"
	"os"
	"strings"

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
		&model.DemandSupport{},
		&model.DemandEvent{},
		&model.DemandComment{},
		&model.DemandCommentReport{},
		&model.Notification{},
	)
}

func backfillOfficeSlugs(db *gorm.DB) error {
	var offices []model.Office
	if err := db.Preload("Councillor.User").Find(&offices).Error; err != nil {
		return err
	}
	for _, office := range offices {
		if strings.TrimSpace(office.Slug) != "" {
			continue
		}
		party, name := "", ""
		if office.Councillor != nil {
			party = office.Councillor.Party
			if office.Councillor.User != nil {
				name = office.Councillor.User.Name
			}
		}
		if err := db.Model(&office).Update("slug", entity.OfficeSlug(party, name, office.CouncillorID)).Error; err != nil {
			return err
		}
	}
	return nil
}

func Run() error {
	openAPISpec, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		return fmt.Errorf("read OpenAPI specification: %w", err)
	}
	if err := environment.Load(); err != nil {
		return err
	}
	mediaStorage, err := platformmedia.NewStorage(environment.Env.Media)
	if err != nil {
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
	if err := backfillOfficeSlugs(db); err != nil {
		return fmt.Errorf("backfill office slugs: %w", err)
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
	// Keep the migration audit event internal while ensuring legacy demands have
	// the same public creation event used by the current timeline and author UI.
	var demandsWithoutCreation []model.Demand
	if err := db.Where("NOT EXISTS (SELECT 1 FROM demand_events WHERE demand_events.demand_id = demands.id AND demand_events.type = ?)", "created").Find(&demandsWithoutCreation).Error; err != nil {
		return err
	}
	for _, demand := range demandsWithoutCreation {
		citizenID := demand.CitizenID
		if err := db.Create(&model.DemandEvent{
			DemandID: demand.ID, Type: "created", ActorUserID: &citizenID,
			Metadata: []byte(`{"source":"legacy_migration"}`), Model: gorm.Model{CreatedAt: demand.CreatedAt},
		}).Error; err != nil {
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

	mediaService := usecase.NewMediaService(mediaStorage)
	authService := usecase.NewAuthService(userRepository, sessionRepository, citizenRepository, councillorRepository, officeMemberRepository, officeMemberRequestRepository, transactionManager, jwtProvider, hashProvider, cryptoProvider, nominatimAddressEnricher, logger)
	authHandler := handler.NewAuthHandler(authService, mediaService)

	officeService := usecase.NewOfficeService(officeRepository, councillorRepository, userRepository, officeMemberRepository, officeMemberRequestRepository, hashProvider, provider.NewSMTPMailer(logger), transactionManager, logger)
	officeHandler := handler.NewOfficeHandler(officeService)
	partyHandler := handler.NewPartyHandler()
	demandService := usecase.NewDemandService(demandRepository, officeRepository, db, logger)
	demandHandler := handler.NewDemandHandler(demandService, officeService, mediaService)

	router := gin.New()
	router.Use(http.ErrorHandler())
	if strings.EqualFold(environment.Env.Media.Driver, "local") {
		router.StaticFS("/media", netHttp.Dir(environment.Env.Media.LocalDir))
	}

	routes.RegisterOpenAPI(router, openAPISpec)
	routes.RegisterAuth(router, authMiddleware, authHandler)
	routes.RegisterParties(router, partyHandler)
	routes.RegisterOffice(router, authMiddleware, officeHandler)
	routes.RegisterDemands(router, authMiddleware, demandHandler)
	routes.RegisterNotifications(router, authMiddleware, demandHandler)

	if err := netHttp.ListenAndServe(":8080", router); err != nil {
		return err
	}

	return nil
}
