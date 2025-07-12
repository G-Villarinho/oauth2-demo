package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type AuthService interface {
	SendVerificationCode(ctx context.Context, email string) (*models.SendVerificationCodeResponse, error)
	Authenticate(ctx context.Context, code, otpID string) (*models.AuthenticateResponse, error)
	ResendVerificationCode(ctx context.Context, otpID string) error
}

type authService struct {
	userRepo   repositories.UserRepository
	otpService OTPService
	jwtService JWTService
	config     configs.Environment
}

func NewAuthService(userRepo repositories.UserRepository, otpService OTPService, jwtService JWTService, config configs.Environment) AuthService {
	return &authService{
		userRepo:   userRepo,
		otpService: otpService,
		jwtService: jwtService,
		config:     config,
	}
}

func (s *authService) SendVerificationCode(ctx context.Context, email string) (*models.SendVerificationCodeResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	otp, err := s.otpService.CreateOTP(ctx, user.ID.Hex(), email)
	if err != nil {
		return nil, fmt.Errorf("create otp: %w", err)
	}

	token, err := s.jwtService.GenerateOTPTokenJWT(ctx, otp.ID.Hex(), otp.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("generate otp token jwt: %w", err)
	}

	// TODO: Send email with token

	return &models.SendVerificationCodeResponse{
		OTPToken:  token,
		ExpiresAt: otp.ExpiresAt,
	}, nil
}

func (s *authService) Authenticate(ctx context.Context, code, otpID string) (*models.AuthenticateResponse, error) {
	otp, err := s.otpService.ValidateCode(ctx, code, otpID)
	if err != nil {
		return nil, fmt.Errorf("validate otp: %w", err)
	}

	expiresAt := time.Now().Add(s.config.OTP.JWTExpirationMinutes)

	token, err := s.jwtService.GenerateAccessTokenJWT(ctx, otp.UserID.Hex(), expiresAt)
	if err != nil {
		return nil, fmt.Errorf("generate access token jwt: %w", err)
	}

	return &models.AuthenticateResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *authService) ResendVerificationCode(ctx context.Context, otpID string) error {
	_, err := s.otpService.ResendCode(ctx, otpID)
	if err != nil {
		return fmt.Errorf("resend otp: %w", err)
	}

	// TODO: Send email with new token

	return nil
}
