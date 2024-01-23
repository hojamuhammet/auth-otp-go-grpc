package otp

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	pb "auth-otp-go-grpc/gen"
	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/internal/database"
	"auth-otp-go-grpc/internal/rabbitmq"
	my_smpp "auth-otp-go-grpc/internal/smpp"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type OTPService struct {
	cfg             *config.Config
	db              *database.Database
	rabbitMQService *rabbitmq.RabbitMQService
	smppConnection  *my_smpp.SMPPConnection
	pb.UnimplementedUserServiceServer
}

func NewOTPService(cfg *config.Config, db *database.Database, rabbitMQService *rabbitmq.RabbitMQService, smppConnection *my_smpp.SMPPConnection) *OTPService {
	return &OTPService{
		cfg:             cfg,
		db:              db,
		rabbitMQService: rabbitMQService,
		smppConnection:  smppConnection,
	}
}

func (s *OTPService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.Empty, error) {
	phoneNumber := req.PhoneNumber

	exists, err := s.checkUserExistenceInDatabase(phoneNumber)
	if err != nil {
		log.Printf("Error checking user existence in the database: %v", err)
		return nil, status.Error(codes.Internal, "Failed to register user")
	}

	var otpCode int32
	if exists {
		otp := GenerateOTP()
		otpCode = int32(otp)
		if err := s.updateUserOTP(phoneNumber, int(otp)); err != nil {
			log.Printf("Error updating OTP in the database: %v", err)
			return nil, status.Error(codes.Internal, "Failed to register user")
		}
	} else {
		otp := GenerateOTP()
		otpCode = int32(otp)
		if err := s.storeUserInDatabase(phoneNumber, int(otp)); err != nil {
			log.Printf("Error creating a new user and storing OTP: %v", err)
			return nil, status.Error(codes.Internal, "Failed to register user")
		}
	}

	err = s.sendDataToRabbitMQ(phoneNumber, otpCode)
	if err != nil {
		log.Printf("Failed to send data to RabbitMQ: %v", err)
		return nil, status.Error(codes.Internal, "Failed to send data to RabbitMQ")
	}

	smsMessage := fmt.Sprintf("Your OTP code is: %d", otpCode)
	if err := s.smppConnection.SendSMS(phoneNumber, smsMessage); err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return nil, status.Error(codes.Internal, "Failed to send SMS")
	}

	return &pb.Empty{}, nil
}

func (s *OTPService) sendDataToRabbitMQ(phoneNumber string, otpCode int32) error {
	data := &pb.RegisterUserResponse{
		User: &pb.User{
			PhoneNumber: phoneNumber,
			Otp:         otpCode,
		},
	}

	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	err = s.rabbitMQService.PublishMessage(context.Background(), "otp_queue", dataBytes)
	if err != nil {
		return err
	}

	return nil
}

func GenerateOTP() int {
	rand.Seed(time.Now().Unix())

	otpCode := 100000 + rand.Intn(900000)

	return otpCode
}

func (s *OTPService) checkUserExistenceInDatabase(phoneNumber string) (bool, error) {
	var exists bool
	err := s.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence in the database: %v", err)

		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (s *OTPService) updateUserOTP(phoneNumber string, newOTP int) error {
	sqlStatement := `
        UPDATE users
        SET otp = $1, otp_created_at = NOW()
        WHERE phone_number = $2
    `

	_, err := s.db.DB.Exec(sqlStatement, newOTP, phoneNumber)
	if err != nil {
		log.Printf("Error updating OTP in the database: %v", err)
		return err
	}

	return nil
}

func (s *OTPService) storeUserInDatabase(phoneNumber string, otp int) error {
	sqlStatement := `
        INSERT INTO users (phone_number, otp, otp_created_at)
        VALUES ($1, $2, NOW())
    `

	_, err := s.db.DB.Exec(sqlStatement, phoneNumber, otp)
	if err != nil {
		log.Printf("Error creating a new user and storing OTP in the database: %v", err)
		return err
	}

	return nil
}

func (s *OTPService) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	phoneNumber := req.PhoneNumber
	otpCode := req.Otp

	var otpFromDB int
	var otpCreatedAt time.Time
	err := s.db.DB.QueryRow("SELECT otp, otp_created_at FROM users WHERE phone_number = $1 AND otp_created_at >= NOW() - interval '5 minutes'", phoneNumber).Scan(&otpFromDB, &otpCreatedAt)

	if err != nil {
		log.Printf("Error querying OTP from the database: %v", err)

		response := &pb.VerifyOTPResponse{
			Valid:    false,
			JwtToken: "",
			Message:  "OTP expired",
		}

		return response, nil
	}

	if otpCode == int32(otpFromDB) {
		jwtToken, err := GenerateJWTToken(phoneNumber, s.cfg.JWTSecret)
		if err != nil {
			log.Printf("Error generating JWT token: %v", err)
			return nil, status.Error(codes.Internal, "OTP verification failed")
		}

		response := &pb.VerifyOTPResponse{
			Valid:    true,
			JwtToken: jwtToken,
			Message:  "OTP verification successful",
		}

		return response, nil
	} else {
		response := &pb.VerifyOTPResponse{
			Valid:    false,
			JwtToken: "",
			Message:  "Invalid OTP",
		}
		return response, nil
	}
}

func GenerateJWTToken(phoneNumber string, jwtSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["phone_number"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (1 day)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetJWTSecretEnv() {
	jwtSecret, err := GenerateJWTSecretKey(64)
	if err != nil {
		log.Fatalf("Failed to generate JWT secret: %v", err)
	}

	os.Setenv("JWT_SECRET", jwtSecret)
}

func GenerateJWTSecretKey(keyLength int) (string, error) {
	keyBytes := make([]byte, keyLength)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(keyBytes), nil
}
