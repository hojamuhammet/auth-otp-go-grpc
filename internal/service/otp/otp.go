package otp

import (
	"context"
	"log/slog"
	"strconv"

	pb "auth-otp-go-grpc/gen"
	"auth-otp-go-grpc/internal/config"
	smpp "auth-otp-go-grpc/internal/service/smpp"
	"auth-otp-go-grpc/pkg/utils"

	repository "auth-otp-go-grpc/internal/repository/postgres"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OTPService struct {
	cfg            config.Config
	smppConnection *smpp.SMPPConnection
	repository     repository.PostgresOTPRepository
	pb.UnimplementedUserServiceServer
}

func NewOTPService(cfg config.Config, smppConnection *smpp.SMPPConnection) *OTPService {
	return &OTPService{
		cfg:            cfg,
		smppConnection: smppConnection,
	}
}

func (s *OTPService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.Empty, error) {
	if ctx.Err() != nil {
		return nil, status.Error(codes.Canceled, "Request canceled")
	}

	phoneNumber := req.PhoneNumber

	exists, err := s.repository.CheckUserExistenceInDatabase(phoneNumber)
	if err != nil {
		slog.Error("Error checking user existence in the database: %v", utils.Err(err))
		return nil, status.Error(codes.Internal, "Failed to register user")
	}

	var otpCode int32
	if exists {
		otp := utils.GenerateOTP()
		otpCode = int32(otp)
		if err := s.repository.UpdateUserOTP(phoneNumber, int(otp)); err != nil {
			slog.Error("Error updating OTP in the database: %v", utils.Err(err))
			return nil, status.Error(codes.Internal, "Failed to register user")
		}
	} else {
		otp := utils.GenerateOTP()
		otpCode = int32(otp)
		if err := s.repository.StoreUserInDatabase(phoneNumber, int(otp)); err != nil {
			slog.Error("Error creating a new user and storing OTP: %v", utils.Err(err))
			return nil, status.Error(codes.Internal, "Failed to register user")
		}
	}

	smsMessage := strconv.Itoa(int(otpCode))
	if err := s.smppConnection.SendSMS(s.cfg, phoneNumber, smsMessage); err != nil {
		slog.Error("Failed to send SMS: %v", err)
		return nil, status.Error(codes.Internal, "Failed to send SMS")
	}

	return &pb.Empty{}, nil
}
