package otp

import (
	"math/rand"
	"time"

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

func GenerateOTP() int {
	// Use the current Unix timestamp (in seconds) as a seed for randomness
	rand.Seed(time.Now().Unix())

	// Generate a random 6-digit integer OTP
	otpCode := 100000 + rand.Intn(900000)

	return otpCode
}
