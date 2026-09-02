package main

import (
	"cidadon/internal/application/usecase"
	"cidadon/internal/platform/config"
	"cidadon/internal/platform/database"
	platformmedia "cidadon/internal/platform/media"
	"context"
	"log"
	"time"

	"go.uber.org/zap"
)

func main() {
	if err := environment.Load(); err != nil {
		log.Fatal(err)
	}
	db, err := database.NewPostgresConnection(environment.Env.Database)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := platformmedia.NewStorage(environment.Env.Media)
	if err != nil {
		log.Fatal(err)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	mediaMigration := usecase.NewMediaMigrationService(db, usecase.NewMediaService(storage), logger.Sugar())
	lifecycle := usecase.NewDemandLifecycleService(db)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		if err := mediaMigration.MigrateLegacyMedia(context.Background()); err != nil {
			log.Printf("migrate legacy media: %v", err)
		}
		if err := lifecycle.CompleteExpiredConfirmations(context.Background(), time.Now()); err != nil {
			log.Printf("complete expired confirmations: %v", err)
		}
		<-ticker.C
	}
}
