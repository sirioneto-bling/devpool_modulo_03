package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/bling-lwsa/devpool-base-web-api/internal/infrastructure/config"
)

// NewConnection opens a MySQL connection using the standard library database/sql.
//
// The go-sql-driver/mysql package is imported with a blank identifier (_) because
// it registers itself via an init() function -- a common Go pattern for database drivers.
// After that, database/sql knows how to speak MySQL.
//
// This function also configures the connection pool and validates the connection
// with a Ping before returning.
func NewConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
