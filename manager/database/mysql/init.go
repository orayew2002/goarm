package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	db *sql.DB
}

func ConnectMySQLFromConfig(cfg Config) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &Connection{db: db}, nil
}

func (c Connection) Exec(ctx context.Context, query string, args ...any) (*result, error) {
	r, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return &result{}, err
	}

	af, err := r.RowsAffected()
	if err != nil {
		return &result{}, err
	}

	return &result{Count: af}, nil
}

func (c Connection) Query(ctx context.Context, query string, args ...any) (*sqlRows, error) {
	r, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return &sqlRows{}, err
	}

	return &sqlRows{rows: r}, nil
}

func (c Connection) Close() error {
	return nil
}

type result struct {
	Count int64
}

func (r *result) RowsAffected() int64 {
	return r.Count
}

type sqlRows struct {
	rows *sql.Rows
}

func (r *sqlRows) Next() bool {
	return r.rows.Next()
}

func (r *sqlRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *sqlRows) Close() {
	_ = r.rows.Close()
}

func (r *sqlRows) Err() error {
	return r.rows.Err()
}
