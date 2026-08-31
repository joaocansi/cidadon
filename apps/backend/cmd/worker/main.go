package main

import (
	"cidadon/internal/application/usecase"
	"cidadon/internal/platform/config"
	"cidadon/internal/platform/database"
	"context"
	"log"
	"time"
)

func main() {
	if err := environment.Load(); err != nil {
		log.Fatal(err)
	}
	db, err := database.NewPostgresConnection(environment.Env.Database)
	if err != nil {
		log.Fatal(err)
	}
	lifecycle := usecase.NewDemandLifecycleService(db)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		if err := lifecycle.CompleteExpiredConfirmations(context.Background(), time.Now()); err != nil {
			log.Printf("complete expired confirmations: %v", err)
		}
		<-ticker.C
	}
}
