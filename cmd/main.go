package main

import (
	"context"
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

	dbInstance, err := database.InitDB(&cfg)
	if err != nil {
		log.Error("Failed to connect to the database: %v", err)
	}
	defer dbInstance.Close()

	grpcServer := server.NewServer(cfg, dbInstance)

	go func() {
		if err := grpcServer.Start(context.Background()); err != nil {
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
