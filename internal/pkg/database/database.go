package database

import (
	"database/sql"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

// Database holds the database connection.
type Database struct {
    db *sql.DB
}

// NewDatabase creates a new Database instance and establishes a database connection.
func NewDatabase(dbURL string) (*Database, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }

    return &Database{db}, nil
}

// Close closes the database connection.
func (d *Database) Close() {
    if d.db != nil {
        d.db.Close()
    }
}
