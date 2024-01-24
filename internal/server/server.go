package server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"sync"

	pb "auth-otp-go-grpc/gen"
	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/internal/database"
	"auth-otp-go-grpc/internal/service/otp"
	smpp "auth-otp-go-grpc/internal/service/smpp"
	"auth-otp-go-grpc/pkg/utils"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg            config.Config
	server         *grpc.Server
	db             *sql.DB
	smppConnection *smpp.SMPPConnection
	pb.UnimplementedUserServiceServer
}

func NewServer(cfg config.Config, dbInstance *database.Database) *Server {
	smppConnection, err := smpp.NewSMPPConnection(cfg) // Initialize and attach SMPP client connection to main listening port for simplicity
	if err != nil {
		slog.Error("Failed to initialize SMPP client: %v", utils.Err(err))
		return nil
	}

	return &Server{
		cfg:            cfg,
		db:             dbInstance.GetDB(),
		smppConnection: smppConnection,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.GrpcServer.Address)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	otpService := otp.NewOTPService(s.cfg, s.smppConnection)
	pb.RegisterUserServiceServer(s.server, otpService)

	reflection.Register(s.server)

	log.Printf("gRPC server started: %s", s.cfg.GrpcServer.Address)

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
