package pgxpool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	db *pgxpool.Pool
}

func ConnectPgxPoolFromConfig(cfg Config) (*Connection, error) {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		sslmode,
	)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return &Connection{db: pool}, nil
}

func (c Connection) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return c.db.Exec(ctx, query, args...)
}

func (c Connection) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return c.db.Query(ctx, query, args...)
}

func (c Connection) Close() error {
	return nil
}
