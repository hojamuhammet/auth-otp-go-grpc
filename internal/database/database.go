package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database holds the database connection.
type Database struct {
    DB *sql.DB
}

// NewDatabase creates a new Database instance and establishes a database connection.
func NewDatabase(dbURL string) (*Database, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    return &Database{db}, nil
}

// Close closes the database connection.
func (d *Database) Close() {
    if d.DB != nil {
        d.DB.Close()
    }
}
