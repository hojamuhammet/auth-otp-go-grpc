package config

import (
	"os"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	GRPCPort     string
	RabbitMQ_URL string
	JWTSecret    string
}

func LoadConfig() Config {
	return Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		GRPCPort:     os.Getenv("GRPC_PORT"),
		RabbitMQ_URL: os.Getenv("RABBITMQ_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
}
