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

// ParseDbTypeFromLabel returns the core database key (e.g., "mysql", "sqlite", "pgxpool")
// from a human-readable label like "MySql" or "Postgres (pgxpool)".
// If the label is unknown, it returns an empty string.
func ParseDbTypeFromLabel(label string) string {
	dbType := DbType(label)
	return dbType.ToCoreDatabase()
}
