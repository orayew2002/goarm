package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

// NewClient creates a new MySQL client with given config.
// If environment variable DB_HOST is set, it overrides cfg.Host.
// Panics on any error.
func NewClient(cfg Config) *sql.DB {
	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		cfg.Host = envHost
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("error opening mysql connection: %v", err))
	}

	// Test connection with Ping
	if err := db.Ping(); err != nil {
		db.Close()
		panic(fmt.Sprintf("can't ping mysql: %v", err))
	}

	return db
}
