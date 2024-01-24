package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/internal/database"
	server "auth-otp-go-grpc/internal/server"
	"auth-otp-go-grpc/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.SetupLogger(cfg.Env)

	slog.Info("Starting the server...", slog.String("env", cfg.Env))
	slog.Debug("Debug messages are enabled") // If env is set to prod, debug messages are going to be disabled

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBname)

	db, err := database.NewDatabase(dbURL)
	if err != nil {
		log.Error("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	grpcServer := server.NewServer(cfg, db)

	go func() {
		if err := grpcServer.Start(context.Background(), cfg); err != nil {
			log.Error("Failed to start gRPC server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Info("Received termination signal. Shutting down...")
	grpcServer.Stop()
	grpcServer.Wait()
}
