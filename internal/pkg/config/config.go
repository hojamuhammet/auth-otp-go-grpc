// config.go
package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    GRPCPort   string
}

// LoadConfig loads the application configuration from environment variables.
func LoadConfig() Config {
    return Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        GRPCPort:   os.Getenv("GRPC_PORT"),
    }
}
