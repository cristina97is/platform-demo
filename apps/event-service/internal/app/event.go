package app

import "time"

// Event описывает событие, которое приходит от клиента.
type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
