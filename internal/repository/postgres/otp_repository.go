package repository

import (
	pb "auth-otp-go-grpc/gen"
	"auth-otp-go-grpc/internal/config"
	"auth-otp-go-grpc/pkg/utils"
	"context"
	"database/sql"
	"time"

	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostgresOTPRepository struct {
	db  *sql.DB
	cfg config.Config
}

func (r *PostgresOTPRepository) CheckUserExistenceInDatabase(phoneNumber string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)", phoneNumber).Scan(&exists)
	if err != nil {
		slog.Info("Error checking user existence in the database: %v", utils.Err(err))

		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

func (r *PostgresOTPRepository) UpdateUserOTP(phoneNumber string, newOTP int) error {
	sqlStatement := `
        UPDATE users
        SET otp = $1, otp_created_at = NOW()
        WHERE phone_number = $2
    `

	_, err := r.db.Exec(sqlStatement, newOTP, phoneNumber)
	if err != nil {
		slog.Error("Error updating OTP in the database: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresOTPRepository) StoreUserInDatabase(phoneNumber string, otp int) error {
	sqlStatement := `
        INSERT INTO users (phone_number, otp, otp_created_at)
        VALUES ($1, $2, NOW())
    `

	_, err := r.db.Exec(sqlStatement, phoneNumber, otp)
	if err != nil {
		slog.Error("Error creating a new user and storing OTP in the database: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresOTPRepository) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	phoneNumber := req.PhoneNumber
	otpCode := req.Otp

	var otpFromDB int
	var otpCreatedAt time.Time
	err := r.db.QueryRow("SELECT otp, otp_created_at FROM users WHERE phone_number = $1 AND otp_created_at >= NOW() - interval '5 minutes'", phoneNumber).Scan(
		&otpFromDB,
		&otpCreatedAt,
	)

	if err != nil {
		slog.Error("Error querying OTP from the database: %v", utils.Err(err))

		response := &pb.VerifyOTPResponse{
			Valid:    false,
			JwtToken: "",
			Message:  "OTP expired",
		}

		return response, nil
	}

	if otpCode == int32(otpFromDB) {
		jwtToken, err := utils.GenerateJWTToken(phoneNumber, r.cfg.JWT.AccessSecretKey)
		if err != nil {
			slog.Error("Error generating JWT token: %v", utils.Err(err))
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
