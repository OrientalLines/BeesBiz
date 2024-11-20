package database

import (
	"embed"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrationFiles defines the order of SQL migration files to be executed
var migrationOrder = []string{
	"definition.sql", // Create tables and basic structure
	"indexes.sql",    // Create indexes for performance
	"ddl.sql",        // Add constraints
	"functions.sql",  // Create functions and procedures
	"trigger.sql",    // Create triggers
}

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB
}

// New creates a new database connection
func New(dataSourceName string) (*DB, error) {
	sqlxDB, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection
	if err := sqlxDB.Ping(); err != nil {
		sqlxDB.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	zap.L().Info("Successfully connected to database")
	return &DB{DB: sqlxDB}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// InitSchema initializes the database schema using embedded SQL files
func (db *DB) InitSchema() error {
	for _, filename := range migrationOrder {
		zap.L().Info("Executing migration", zap.String("file", filename))

		// Read the SQL file from embedded FS
		content, err := migrations.ReadFile("migrations/" + filename)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", filename, err)
		}

		// Execute the SQL
		if err := db.ExecuteSQL(string(content)); err != nil {
			zap.L().Error("Failed to execute migration",
				zap.String("file", filename),
				zap.Error(err))
			return fmt.Errorf("error executing migration %s: %w", filename, err)
		}

		zap.L().Info("Successfully executed migration", zap.String("file", filename))
	}

	zap.L().Info("All migrations completed successfully")
	return nil
}

// ExecuteSQL executes a SQL string
func (db *DB) ExecuteSQL(sql string) error {
	// Split the SQL into individual statements (if multiple are present)
	statements := strings.Split(sql, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("error executing SQL statement: %w", err)
		}
	}

	return nil
}
