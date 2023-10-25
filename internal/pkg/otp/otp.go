package otp

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/hojamuhammet/go-grpc-otp-rabbitmq/gen"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/config"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/database"
	"github.com/hojamuhammet/go-grpc-otp-rabbitmq/internal/pkg/rabbitmq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type OTPService struct {
	cfg *config.Config
	db *database.Database
    rabbitMQService *rabbitmq.RabbitMQService
	pb.UnimplementedUserServiceServer
}

func NewOTPService(cfg *config.Config, db *database.Database, rabbitMQService *rabbitmq.RabbitMQService) *OTPService {
	return &OTPService{
		cfg: cfg,
		db: db,
        rabbitMQService: rabbitMQService,
	}
}

// RegisterUser handles user registration and OTP generation.
func (s *OTPService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
    phoneNumber := req.PhoneNumber

    // Check if the user with the given phone number exists
    exists, err := s.checkUserExistenceInDatabase(phoneNumber)
    if err != nil {
        log.Printf("Error checking user existence in the database: %v", err)
        return nil, status.Error(codes.Internal, "Failed to register user")
    }

    var otpCode int32
    if exists {
        // User exists, update their OTP in the database
        otp := GenerateOTP()
        otpCode = int32(otp)
        if err := s.updateUserOTP(phoneNumber, int(otp)); err != nil {
            log.Printf("Error updating OTP in the database: %v", err)
            return nil, status.Error(codes.Internal, "Failed to register user")
        }
    } else {
        // User doesn't exist, create a new user and store the OTP
        otp := GenerateOTP()
        otpCode = int32(otp)
        if err := s.storeUserInDatabase(phoneNumber, int(otp)); err != nil {
            log.Printf("Error creating a new user and storing OTP: %v", err)
            return nil, status.Error(codes.Internal, "Failed to register user")
        }
    }

    // Prepare the response message
    response := &pb.RegisterUserResponse{
        User: &pb.User{
            PhoneNumber: phoneNumber,
            Otp:         otpCode,
        },
    }

    log.Printf("Registered user with phone number: %s", phoneNumber)

    // Marshal the response into Protocol Buffers binary format
    responseBytes, err := proto.Marshal(response)
    if err != nil {
        log.Printf("Error marshaling response: %v", err)
        return nil, status.Error(codes.Internal, "Failed to register user")
    }

    // Send the response to RabbitMQ using the PublishMessage function
    err = s.rabbitMQService.PublishMessage(ctx, "otp_queue", responseBytes)
    if err != nil {
        log.Printf("Failed to send message to RabbitMQ: %v", err)
        // Handle the error or return an appropriate gRPC error.
    }

    // Return the response
    return response, nil
}

func GenerateOTP() int {
    // Use the current Unix timestamp (in seconds) as a seed for randomness
    rand.Seed(time.Now().Unix())

    // Generate a random 6-digit integer OTP
    otpCode := 100000 + rand.Intn(900000)

    return otpCode
}

// GenerateJWTToken generates a JWT token for the given phone number.
func GenerateJWTToken(phoneNumber string, secretKey string) (string, error) {
	// Define the claims for the JWT token
	claims := jwt.MapClaims{
		"phone_number": phoneNumber,
		"exp":          time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *OTPService) checkUserExistenceInDatabase(phoneNumber string) (bool, error) {
    // Query the database to check if a user with the given phone number exists
    var exists bool
    err := s.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
    if err != nil {
        log.Printf("Error checking user existence in the database: %v", err)

        // Check if the error is due to a database connection issue
        if err == sql.ErrNoRows {
            return false, nil
        }

        return false, err
    }

    return exists, nil
}

// updateUserOTP updates the OTP for a user in the database.
func (s *OTPService) updateUserOTP(phoneNumber string, newOTP int) error {
    // Prepare the SQL statement to update the OTP for the user with the given phone number
    sqlStatement := `
        UPDATE users
        SET otp = $1, otp_created_at = NOW()
        WHERE phone_number = $2
    `

    // Execute the SQL statement
    _, err := s.db.DB.Exec(sqlStatement, newOTP, phoneNumber)
    if err != nil {
        log.Printf("Error updating OTP in the database: %v", err)
        return err
    }

    return nil
}

// storeUserInDatabase creates a new user in the database with the given phone number and OTP,
func (s *OTPService) storeUserInDatabase(phoneNumber string, otp int) error {
    sqlStatement := `
        INSERT INTO users (phone_number, otp, otp_created_at)
        VALUES ($1, $2, CURRENT_TIMESTAMP)
    `

    // Execute the SQL statement
    _, err := s.db.DB.Exec(sqlStatement, phoneNumber, otp)
    if err != nil {
        log.Printf("Error creating a new user and storing OTP in the database: %v", err)
        return err
    }

    return nil
}
