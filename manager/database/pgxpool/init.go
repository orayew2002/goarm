package pgxpool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

// NewClient creates a new PostgreSQL connection pool and panics on failure.
func NewClient(cfg Config) *pgxpool.Pool {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		sslmode,
	)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(fmt.Sprintf("pgxpool.New error: %v", err))
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		panic(fmt.Sprintf("PostgreSQL ping failed: %v", err))
	}

	return pool
}
