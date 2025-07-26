package manager

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed database/mysql/*
var databaseMysql embed.FS

//go:embed database/pgxpool/*
var databasePsql embed.FS

//go:embed database/sqlite/*
var databaseSqlite embed.FS

// driverFS maps database type to its embedded FS
var driverFS = map[string]embed.FS{
	"sqlite":  databaseSqlite,
	"pgxpool": databasePsql,
	"mysql":   databaseMysql,
}

// Manager holds embedded database files
type Manager struct {
	Database DatabaseFiles
}

// DatabaseFiles holds the key source files for a database driver
type DatabaseFiles struct {
	config []byte
	domain []byte
	init   []byte
}

func (df *DatabaseFiles) GetConfig() []byte {
	return df.config
}

func (df *DatabaseFiles) GetDomain() []byte {
	return df.domain
}

func (df *DatabaseFiles) GetInit() []byte {
	return df.init
}

// Manage loads and returns all embedded files for the given database type
func Manage(database string) *Manager {
	fsys, ok := driverFS[database]
	if !ok {
		panic(fmt.Errorf("unknown database driver: %s", database))
	}

	basePath := fmt.Sprintf("database/%s", database)

	config, err := readFile(fsys, basePath+"/config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load config.yaml: %w", err))
	}

	domain, err := readFile(fsys, basePath+"/domain.go")
	if err != nil {
		panic(fmt.Errorf("failed to load domain.go: %w", err))
	}

	initCode, err := readFile(fsys, basePath+"/init.go")
	if err != nil {
		panic(fmt.Errorf("failed to load init.go: %w", err))
	}

	return &Manager{
		Database: DatabaseFiles{
			config: config,
			domain: domain,
			init:   initCode,
		},
	}
}

// readFile is a small helper to read from embed.FS with a clear error
func readFile(fsys fs.FS, path string) ([]byte, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %w", path, err)
	}
	return data, nil
}
