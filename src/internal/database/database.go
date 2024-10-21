package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB
}

var db *DB

// New creates a new database connection
func New(dataSourceName string) error {
	sqlxDB, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection
	if err := sqlxDB.Ping(); err != nil {
		sqlxDB.Close() // Close the connection if ping fails
		return fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to database")
	db = &DB{DB: sqlxDB}

	return nil
}

// GetDB returns the singleton database instance
func GetDB() *DB {
	return db
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

func (db *DB) InitSchema() error {
	return nil
}
