package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrationFiles defines the order of SQL migration files to be executed
var migrationOrder = []string{
	"definition.sql", // Create tables and basic structure
	"ddl.sql",        // Add constraints
	"indexes.sql",    // Create indexes for performance
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
		if _, err := db.Exec(string(content)); err != nil {
			// Get the detailed error message
			sqlErr, ok := err.(*pq.Error)
			if ok {
				zap.L().Error("Failed to execute migration",
					zap.String("file", filename),
					zap.String("error_detail", sqlErr.Detail),
					zap.String("error_hint", sqlErr.Hint),
					zap.String("error_position", sqlErr.Position),
					zap.String("error_where", sqlErr.Where),
					zap.String("error_schema", sqlErr.Schema),
					zap.String("error_table", sqlErr.Table),
					zap.String("error_column", sqlErr.Column),
					zap.String("error_data_type", sqlErr.DataTypeName),
					zap.String("error_constraint", sqlErr.Constraint),
					zap.Error(err))
				return fmt.Errorf("error executing migration %s at position %s: %w", filename, sqlErr.Position, err)
			}
			
			// Fallback for non-postgres errors
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
