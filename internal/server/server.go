package server

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/config"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/database"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/otp"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/rabbitmq"
	my_smpp "github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/smpp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg             *config.Config
	server          *grpc.Server
	db              *database.Database
	rabbitMQService *rabbitmq.RabbitMQService
	smppConnection  *my_smpp.SMPPConnection
	pb.UnimplementedUserServiceServer
}

func NewServer(cfg *config.Config, db *database.Database, rabbitMQService *rabbitmq.RabbitMQService) *Server {
	smppConnection, err := my_smpp.NewSMPPConnection() // Initialize the SMPP client
	if err != nil {
		log.Fatalf("Failed to initialize SMPP client: %v", err)
		return nil
	}

	return &Server{
		cfg:             cfg,
		db:              db,
		rabbitMQService: rabbitMQService,
		smppConnection:  smppConnection,
	}
}

func (s *Server) Start(ctx context.Context, cfg *config.Config) error {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	otpService := otp.NewOTPService(s.cfg, s.db, s.rabbitMQService, s.smppConnection)
	pb.RegisterUserServiceServer(s.server, otpService)

	reflection.Register(s.server)

	log.Printf("gRPC server started on port %s", s.cfg.GRPCPort)

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

func (s *Server) Wait() {
	var wg sync.WaitGroup
	wg.Wait()
}
