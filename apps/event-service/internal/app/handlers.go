package app

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// handleHealth показывает, что процесс приложения жив.
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// handleReady показывает, что приложение готово работать.
// Здесь мы уже проверяем доступность базы данных.
func (a *App) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.storage.Ping(ctx); err != nil {
		http.Error(w, "database is not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ready"))
}

// handleEvents маршрутизирует методы GET и POST.
func (a *App) handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.createEvent(w, r)
	case http.MethodGet:
		a.listEvents(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// createEvent принимает JSON и сохраняет событие в PostgreSQL.
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

// listEvents читает события из PostgreSQL и возвращает их клиенту.
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
