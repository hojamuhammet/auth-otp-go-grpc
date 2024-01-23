package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/internal/database"
	"auth-otp-go-grpc/internal/rabbitmq"
	server "auth-otp-go-grpc/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the env variables: %v", err)
	}

	cfg := config.LoadConfig()

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := database.NewDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	rabbitMQService, err := rabbitmq.InitRabbitMQConnection(cfg.RabbitMQ_URL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitMQService.Close()

	grpcServer := server.NewServer(&cfg, db, rabbitMQService)

	go func() {
		if err := grpcServer.Start(context.Background(), &cfg); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Received termination signal. Shutting down...")
	grpcServer.Stop()
	grpcServer.Wait()
}
