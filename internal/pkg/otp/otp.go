package otp

import (
	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/config"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/database"
)

type OTPService struct {
	cfg *config.Config
	db *database.Database
	pb.UnimplementedUserServiceServer
}

func NewOTPService(cfg *config.Config, db *database.Database) *OTPService {
	return &OTPService{
		cfg: cfg,
		db: db,
	}
}
