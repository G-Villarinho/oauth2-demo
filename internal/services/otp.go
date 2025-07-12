package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	otpLength      = 6
	otpChars       = "0123456789"
	ResendCooldown = 60 * time.Second
)

type OTPService interface {
	CreateOTP(ctx context.Context, userID, email string) (*entities.OTP, error)
	ValidateCode(ctx context.Context, code, otpID string) (*entities.OTP, error)
	ResendCode(ctx context.Context, otpID string) (*entities.OTP, error)
}

type otpService struct {
	otpRepo repositories.OTPRepository
	config  *configs.Environment
}

func NewOTPService(otpRepo repositories.OTPRepository, config *configs.Environment) OTPService {
	return &otpService{
		otpRepo: otpRepo,
		config:  config,
	}
}

func (s *otpService) CreateOTP(ctx context.Context, userID, email string) (*entities.OTP, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("convert userID to ObjectID: %w", err)
	}

	otp := &entities.OTP{
		UserID:    userIDObj,
		Email:     email,
		Code:      s.generateOTP(),
		ExpiresAt: time.Now().UTC().Add(s.config.OTP.ExpirationMinutes),
	}

	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("create otp: %w", err)
	}

	return otp, nil
}

func (s *otpService) ValidateCode(ctx context.Context, code, otpID string) (*entities.OTP, error) {
	otp, err := s.otpRepo.FindByID(ctx, otpID)
	if err != nil {
		return nil, fmt.Errorf("find otp by id: %w", err)
	}

	if err := otp.ValidateCode(code); err != nil {
		return nil, fmt.Errorf("validate code: %w", err)
	}

	if err := s.otpRepo.Delete(ctx, otpID); err != nil {
		return nil, fmt.Errorf("delete otp: %w", err)
	}

	return otp, nil
}

func (s *otpService) ResendCode(ctx context.Context, otpID string) (*entities.OTP, error) {
	otp, err := s.otpRepo.FindByID(ctx, otpID)
	if err != nil {
		return nil, fmt.Errorf("find otp by id: %w", err)
	}

	if !otp.IsResendable() {
		return nil, &domain.ErrOTPNotResendable{
			TimeRemaining: ResendCooldown - time.Since(*otp.ResendAt),
		}
	}

	otp.Code = s.generateOTP()

	if err := s.otpRepo.UpdateCode(ctx, otpID, otp.Code); err != nil {
		return nil, fmt.Errorf("update otp code: %w", err)
	}

	return otp, nil
}

func (s *otpService) generateOTP() string {
	otp := make([]byte, otpLength)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}
	return string(otp)
}
