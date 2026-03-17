package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cristina97is/platform-demo/apps/event-service/internal/config"
)

// EventRow описывает запись события в базе данных.
type EventRow struct {
	ID        int
	UserID    int
	Type      string
	Amount    int
	CreatedAt time.Time
}

// PostgresStorage отвечает за работу с PostgreSQL.
type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage создаёт подключение к PostgreSQL.
func NewPostgresStorage(cfg config.Config) (*PostgresStorage, error) {
	// Формируем DSN — строку подключения к базе.
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Создаём пул соединений с базой.
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		pool: pool,
	}, nil
}

// Close закрывает соединения с базой.
func (s *PostgresStorage) Close() {
	s.pool.Close()
}

// Ping проверяет, доступна ли база данных.
func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

// CreateEvent сохраняет событие в базу и возвращает созданную запись.
func (s *PostgresStorage) CreateEvent(ctx context.Context, userID int, eventType string, amount int) (EventRow, error) {
	query := `
		INSERT INTO events (user_id, type, amount)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, type, amount, created_at
	`

	var row EventRow

	err := s.pool.QueryRow(ctx, query, userID, eventType, amount).Scan(
		&row.ID,
		&row.UserID,
		&row.Type,
		&row.Amount,
		&row.CreatedAt,
	)
	if err != nil {
		return EventRow{}, err
	}

	return row, nil
}

// ListEvents возвращает все события из базы по убыванию ID.
func (s *PostgresStorage) ListEvents(ctx context.Context) ([]EventRow, error) {
	query := `
		SELECT id, user_id, type, amount, created_at
		FROM events
		ORDER BY id DESC
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]EventRow, 0)

	for rows.Next() {
		var row EventRow

		err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.Type,
			&row.Amount,
			&row.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, rows.Err()
}
