package domain

// App holds metadata about the application.
type App struct {
	Name   string
	DbType DbType
}

// DbType represents a supported database type.
type DbType string

const (
	DBTypePostgres DbType = "Postgres (pgxpool)"
	DBTypeMySQL    DbType = "MySql"
	DBTypeSQLite   DbType = "Sqlite"
)

// SupportedDatabaseTypes lists all available database types.
var SupportedDatabaseTypes = []DbType{
	DBTypePostgres,
	DBTypeMySQL,
	DBTypeSQLite,
}

// ToCoreDatabase returns the internal key used for this DbType (e.g., for folder or driver names).
func (d DbType) ToCoreDatabase() string {
	switch d {
	case DBTypeMySQL:
		return "mysql"
	case DBTypeSQLite:
		return "sqlite"
	case DBTypePostgres:
		return "pgxpool"
	default:
		return ""
	}
}

// ToCoreDatabase returns the internal key used for this DbType (e.g., for folder or driver names).
func (d DbType) ToCoreConfig() string {
	switch d {
	case DBTypeMySQL:
		return "Mysql"
	case DBTypeSQLite:
		return "Sqlite"
	case DBTypePostgres:
		return "Psql"
	default:
		return ""
	}
}

// ToCoreDatabase returns the internal key used for this DbType (e.g., for folder or driver names).
func (d DbType) PackagePath() string {
	switch d {
	case DBTypeMySQL:
		return "database/sql"
	case DBTypeSQLite:
		return "database/sql"
	case DBTypePostgres:
		return "github.com/jackc/pgx/v5/pgxpool"
	default:
		return ""
	}
}

// ToCoreDatabase returns the internal key used for this DbType (e.g., for folder or driver names).
func (d DbType) PackageVal() string {
	switch d {
	case DBTypeMySQL:
		return "sql.DB"
	case DBTypeSQLite:
		return "sql.DB"
	case DBTypePostgres:
		return "pgxpool.Pool"
	default:
		return ""
	}
}

func (d DbType) GetDockerConfig() string {
	switch d {
	case DBTypePostgres:
		return `
  db:
    image: postgres:16
    container_name: postgres_db
    environment:
      POSTGRES_DB: your_database
      POSTGRES_USER: your_username
      POSTGRES_PASSWORD: your_password
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    networks:
      - app-network
`
	case DBTypeMySQL:
		return `
  db:
    image: mysql:8
    container_name: mysql_db
    environment:
      MYSQL_DATABASE: your_database
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_USER: your_username
      MYSQL_PASSWORD: your_password
    ports:
      - "3306:3306"
    volumes:
      - mysqldata:/var/lib/mysql
    networks:
      - app-network
`
	case DBTypeSQLite:
		return `
  # SQLite is file-based and doesn't run as a container.
  # Your application must mount a shared volume to persist the .db file.

  # Example volume mount in app service (not a real DB container):
  # volumes:
  #   - ./data:/app/data
`
	default:
		return ""
	}
}

// ParseDbTypeFromLabel returns the core database key (e.g., "mysql", "sqlite", "pgxpool")
// from a human-readable label like "MySql" or "Postgres (pgxpool)".
// If the label is unknown, it returns an empty string.
func ParseDbTypeFromLabel(label string) string {
	dbType := DbType(label)
	return dbType.ToCoreDatabase()
}
