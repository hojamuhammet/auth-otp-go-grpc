package otp

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
)

type OTPService struct {
	db *sql.DB
}

func NewOTPService(db *sql.DB) *OTPService {
	return &OTPService{
		db: db,
	}
}

// CheckPhoneNumber is the implementation of the CheckPhoneNumber gRPC method.
func (s *OTPService) CheckPhoneNumber(ctx context.Context, req *pb.CheckPhoneNumberRequest) (*pb.CheckPhoneNumberResponse, error) {
	phoneNumber := req.PhoneNumber

	// Query the database to check if the phone number exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check phone number existence: %v", err)
		return nil, err
	}

	// Create and return the response
	response := &pb.CheckPhoneNumberResponse{
		Exists: exists,
	}

	log.Printf("Checked phone number existence for %s: %v", phoneNumber, exists)

	return response, nil
}

