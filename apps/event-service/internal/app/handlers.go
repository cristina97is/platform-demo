package app

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// -----------------------------
// Prometheus метрики
// -----------------------------

var requestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"endpoint", "method"},
)

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests",
	},
	[]string{"endpoint"},
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

// -----------------------------
// Handlers
// -----------------------------

// Главная страница сервиса
func (a *App) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]any{
		"service": "event-service",
		"status":  "running",
		"endpoints": []string{
			"/",
			"/healthz",
			"/readyz",
			"/events",
			"/metrics",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Проверка, жив ли процесс
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Проверка готовности сервиса
func (a *App) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.storage.Ping(ctx); err != nil {
		http.Error(w, "database not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}

// Роутер для /events
func (a *App) handleEvents(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestsTotal.WithLabelValues("/events", r.Method).Inc()

	defer func() {
		requestDuration.
			WithLabelValues("/events").
			Observe(time.Since(start).Seconds())
	}()

	switch r.Method {
	case http.MethodPost:
		a.createEvent(w, r)
	case http.MethodGet:
		a.listEvents(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Создание события
func (a *App) createEvent(w http.ResponseWriter, r *http.Request) {
	var input Event

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.UserID <= 0 {
		http.Error(w, "user_id must be greater than 0", http.StatusBadRequest)
		return
	}

	if input.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	row, err := a.storage.CreateEvent(ctx, input.UserID, input.Type, input.Amount)
	if err != nil {
		http.Error(w, "failed to create event", http.StatusInternalServerError)
		return
	}

	result := Event{
		ID:        row.ID,
		UserID:    row.UserID,
		Type:      row.Type,
		Amount:    row.Amount,
		CreatedAt: row.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
}

// Получение списка событий
func (a *App) listEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := a.storage.ListEvents(ctx)
	if err != nil {
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	result := make([]Event, 0, len(rows))

	for _, row := range rows {
		result = append(result, Event{
			ID:        row.ID,
			UserID:    row.UserID,
			Type:      row.Type,
			Amount:    row.Amount,
			CreatedAt: row.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(result)
}
