package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/zap"
)

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

func (db *DB) InitSchema(pathToScripts string, sqlFiles []string) error {
	for _, file := range sqlFiles {
		filePath := filepath.Join(pathToScripts, file)
		zap.L().Info("Loading SQL file", zap.String("file", file))
		if err := db.executeSQLFile(filePath); err != nil {
			zap.L().Error("Failed to execute SQL file", zap.String("file", file), zap.Error(err))
			return fmt.Errorf("error executing SQL file %s: %w", file, err)
		}
		zap.L().Info("Successfully executed SQL file", zap.String("file", file))
	}

	zap.L().Info("All SQL files executed successfully")
	return nil
}

func (db *DB) executeSQLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error executing SQL: %w", err)
	}
	return nil
}

func (db *DB) ExecuteSQL(sql string) error {
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error executing SQL: %w", err)
	}
	return nil
}
