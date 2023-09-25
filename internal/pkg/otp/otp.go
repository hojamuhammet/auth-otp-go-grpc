package otp

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	utils "github.com/hojamuhammet/go-grpc-otp-rabbitmq/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OTPService struct {
	db *sql.DB
}

func NewOTPService(db *sql.DB) *OTPService {
	return &OTPService {
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

func (s *OTPService) GenerateOTP(ctx context.Context, req *pb.GenerateOTPRequest) (*pb.GenerateOTPResponse, error) {
	// Extract user information from the request
	userID := req.UserId
	phoneNumber := req.PhoneNumber

	// Generate a 6-digit OTP
	otp := utils.GenerateRandomOTP()

	// Update the user's OTP and OTP creation timestamp in the database
	if err := s.updateUserOTP(userID, otp); err != nil {
		log.Printf("Error updating OTP in the database: %v", err)
		return nil, status.Error(codes.Internal, "Failed to generate and save OTP")
	}

	// Create and return the response
	response := &pb.GenerateOTPResponse{
		Otp: otp,
	}

	log.Printf("Generated OTP %s for user %d with phone number %s", otp, userID, phoneNumber)

	return response, nil
}


func (s *OTPService) updateUserOTP(userID int64, otp string) error {
	sqlStatement := `
        UPDATE users
        SET otp = $1, otp_created_at = NOW()
        WHERE id = $2
    `

	_, err := s.db.Exec(sqlStatement, otp, userID)
	if err != nil {
		log.Printf("Error updating OTP in the database: %v", err)
		return err
	}

	return nil
}
