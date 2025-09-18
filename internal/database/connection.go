package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Connect establishes a connection to the SQLite database
func Connect(databaseURL string) (*sql.DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(databaseURL)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Connect to database
	db, err := sql.Open("sqlite3", databaseURL+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Info().Str("database", databaseURL).Msg("Connected to database")
	return db, nil
}