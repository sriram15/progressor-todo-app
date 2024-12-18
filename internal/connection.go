package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/glebarez/sqlite"
	"github.com/pressly/goose/v3"
)

const MIGRATION_PATH = "internal/database/migrations"
const DATABASE_NAME = "progressor.db"

// getDatabasePath returns the path to the SQLite database file.
func GetDatabasePath() (string, error) {

	// Check if DATBASE_PATH env variable is set
	if databasePathEnv := os.Getenv("DATABASE_PATH"); databasePathEnv != "" {
		databasePathEnv = filepath.Join(databasePathEnv, DATABASE_NAME)
		return databasePathEnv, nil
	}

	var appDir string
	// Determine the appropriate directory based on OS
	switch runtime.GOOS {
	case "windows":
		appDir = os.Getenv("APPDATA")
	case "darwin": // macOS
		appDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		appDir = filepath.Join(os.Getenv("HOME"), ".config")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	// Create a directory for your app
	appDir = filepath.Join(appDir, APP_NAME)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", err
	}

	// Define the database file path
	dbPath := filepath.Join(appDir, DATABASE_NAME)
	return dbPath, nil
}

func OpenDB() (*sql.DB, error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %v", err)
	}

	db, err := goose.OpenDBWithDriver("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("goose: failed to open DB: %v", err)
	}

	if err := goose.Up(db, MIGRATION_PATH); err != nil {
		return nil, fmt.Errorf("failed to create migration: %v", err)
	}

	fmt.Println("Migration created successfully.")
	return db, nil
}

func CloseGooseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		fmt.Println("Failed to close DB:", err)
	}
}
