package otp

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OTPService struct {
	db *sql.DB
}

func NewOTPService(db *sql.DB) *OTPService {
	return &OTPService{
		db: db,
	}
}

func (s *OTPService) CheckPhoneNumber(ctx context.Context, req *pb.CheckPhoneNumberRequest) (*pb.CheckPhoneNumberResponse, error) {
	phoneNumber := req.PhoneNumber

	// Query the database to check if the phone number exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	if err != nil {
		log.Printf("Error checking phone number existence in the database: %v", err)

		// Check if the error is due to a database connection issue
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Phone number not found")
		}

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	response := &pb.CheckPhoneNumberResponse{
		Exists: exists,
	}

	log.Printf("Checked phone number existence for %s: %v", phoneNumber, exists)

	return response, nil
}

