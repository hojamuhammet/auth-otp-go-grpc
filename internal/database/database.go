package database

import (
	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/pkg/utils"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq" // init postgresql driver
)

type Database struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) (*Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode)

	dbInstance, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("failed to initialize database: %v", utils.Err(err))
		return nil, err
	}

	if err := dbInstance.Ping(); err != nil {
		slog.Error("failed to initialize database: %v", utils.Err(err))
		return nil, err
	}

	return &Database{db: dbInstance}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
