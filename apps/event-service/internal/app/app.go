package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cristina97is/platform-demo/apps/event-service/internal/config"
	"github.com/cristina97is/platform-demo/apps/event-service/internal/storage"
)

// App — корневой объект приложения.
type App struct {
	cfg     config.Config
	router  *http.ServeMux
	storage *storage.PostgresStorage
}

// New создаёт экземпляр приложения.
func New(cfg config.Config, storage *storage.PostgresStorage) *App {
	mux := http.NewServeMux()

	a := &App{
		cfg:     cfg,
		router:  mux,
		storage: storage,
	}

	a.routes()

	return a
}

// Router возвращает HTTP-роутер приложения.
func (a *App) Router() http.Handler {
	return a.router
}

// routes регистрирует маршруты.
func (a *App) routes() {
	a.router.HandleFunc("/", a.handleRoot)
	a.router.HandleFunc("/healthz", a.handleHealth)
	a.router.HandleFunc("/readyz", a.handleReady)
	a.router.HandleFunc("/events", a.handleEvents)
	a.router.Handle("/metrics", promhttp.Handler())
}
