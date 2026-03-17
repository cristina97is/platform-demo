package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cristina97is/platform-demo/apps/event-service/internal/app"
	"github.com/cristina97is/platform-demo/apps/event-service/internal/config"
	"github.com/cristina97is/platform-demo/apps/event-service/internal/storage"
)

func main() {
	// Загружаем конфигурацию приложения.
	cfg := config.Load()

	// Создаём storage для PostgreSQL.
	pgStorage, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pgStorage.Close()

	// Создаём приложение.
	application := app.New(cfg, pgStorage)

	// Создаём HTTP-сервер.
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           application.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting event-service on port %s", cfg.Port)

	// Запускаем сервер.
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
