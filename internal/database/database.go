package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

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

func (d *Database) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
