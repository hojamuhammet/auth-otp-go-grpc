package server

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/config"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/database"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/otp"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents your gRPC server.
type Server struct {
    cfg *config.Config
    server *grpc.Server
    db *database.Database
    rabbitMQService *rabbitmq.RabbitMQService
    pb.UnimplementedUserServiceServer
}

// NewServer creates a new instance of the Server.
func NewServer(cfg *config.Config, db *database.Database, rabbitMQService *rabbitmq.RabbitMQService) *Server {
    return &Server{
        cfg: cfg,
        db: db,
        rabbitMQService: rabbitMQService,
    }
}

// Start starts the gRPC server.
func (s *Server) Start(ctx context.Context, cfg *config.Config) error {
    lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
    if err != nil {
        return err
    }

    s.server = grpc.NewServer()

    otpService := otp.NewOTPService(s.cfg, s.db, s.rabbitMQService)
    pb.RegisterUserServiceServer(s.server, otpService)

    reflection.Register(s.server)

    log.Printf("gRPC server started on port %s", s.cfg.GRPCPort)

    return s.server.Serve(lis)
}

// Stop stops the gRPC server gracefully.
func (s *Server) Stop() {
    if s.server != nil {
        s.server.GracefulStop()
    }
}

// Wait waits for the server to finish gracefully.
func (s *Server) Wait() {
    var wg sync.WaitGroup
    wg.Wait()
}