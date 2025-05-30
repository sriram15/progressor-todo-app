package connection

import (
	"database/sql"
	"fmt"
	"os"
)

// DATABASE_NAME is the default name for the SQLite database file.
const DATABASE_NAME = "progressor.db"

var (
	connector DBConnector
	db        *sql.DB
)

// OpenDB determines the database type from environment variables,
// creates the appropriate connector, connects, and applies migrations.
func OpenDB() (*sql.DB, error) {
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case DBTypeSQLite, "sqlite3", "": // Default to SQLite
		fmt.Println("Using SQLite database connector.")
		connector = NewSQLiteConnector()
	case DBTypeTurso:
		fmt.Println("Using Turso database connector.")
		connector = NewTursoConnector()
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s. Supported types are 'sqlite' or 'turso'", dbType)
	}

	connectedDB, actualDBType, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", dbType, err)
	}

	fmt.Printf("Successfully established database connection (Type: %s).\n", actualDBType)

	if err := connector.Migrate(connectedDB, actualDBType); err != nil {
		// It's important to close the DB if migration fails to prevent leaks.
		connectedDB.Close()
		return nil, fmt.Errorf("failed to apply migrations for %s: %w", actualDBType, err)
	}

	db = connectedDB

	fmt.Println("Database ready.")
	return db, nil
}

// GetOrReconnectDB checks if the current db connection is alive, and if not, reopens it.
func GetOrReconnectDB() (*sql.DB, error) {
	if db != nil {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		fmt.Println("DB connection lost, attempting to reconnect...")
	}
	// (Re)open the DB connection
	return OpenDB()
}

func GetDBInfo() (string, string) {
	return connector.GetDBInfo()
}
