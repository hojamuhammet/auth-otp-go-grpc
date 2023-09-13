package otp

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// OTPService represents the OTP generation and validation service.
type OTPService struct {
	db *sql.DB
}

// NewOTPService creates a new instance of the OTPService.
func NewOTPService(db *sql.DB) *OTPService {
	return &OTPService{db}
}

// GenerateOTP generates a random 6-digit OTP for the given user ID and phone number.
func (s *OTPService) GenerateOTP(userID int64, phoneNumber string) (string, error) {
	// Generate a 6-digit random OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store the OTP in the database along with the user's phone number
	_, err := s.db.Exec("INSERT INTO users (phone_number) VALUES ($1)", phoneNumber)
	if err != nil {
		log.Printf("Failed to store OTP in the database: %v", err)
		return "", err
	}

	log.Printf("Generated OTP %s for phone number %s", otp, phoneNumber)
	return otp, nil
}

// VerifyOTP validates the OTP for the given user ID and phone number.
func (s *OTPService) VerifyOTP(userID int64, phoneNumber string, otp string) error {
	// Retrieve the stored OTP for the given phone number from the database
	var storedOTP string
	err := s.db.QueryRow("SELECT phone_number FROM users WHERE phone_number = $1", phoneNumber).Scan(&storedOTP)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Phone number not found: %s", phoneNumber)
			return errors.New("Phone number not found")
		}
		log.Printf("Error while querying the database: %v", err)
		return err
	}

	if otp != storedOTP {
		log.Printf("Invalid OTP for phone number %s", phoneNumber)
		return errors.New("Invalid OTP")
	}

	log.Printf("Successfully verified OTP for phone number %s", phoneNumber)

	// You can implement OTP expiration logic here if needed

	return nil
}
