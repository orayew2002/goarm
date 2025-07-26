package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Path string `yaml:"path"`
}

type connection struct {
	db *sql.DB
}

// NewClient creates a SQLite connection.
// If the database file or its directory doesn't exist, they will be created.
// Panics on any failure.
func NewClient(config Config) *sql.DB {
	if config.Path == "" {
		panic("sqlite config error: database path is empty")
	}

	// Ensure parent directory exists
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("sqlite directory creation error: %v", err))
	}

	// Ensure the file exists (create if missing)
	file, err := os.OpenFile(config.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("sqlite file creation/opening error: %v", err))
	}
	file.Close()

	// Open DB connection
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		panic(fmt.Sprintf("sqlite open error: %v", err))
	}

	// Ping to verify
	if err := db.Ping(); err != nil {
		_ = db.Close()
		panic(fmt.Sprintf("sqlite ping error: %v", err))
	}

	return db
}
