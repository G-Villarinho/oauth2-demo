package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	otpLength = 6
	otpChars  = "0123456789"
)

type OTPService interface {
	CreateOTP(ctx context.Context, userID string) (*models.OTP, error)
	ValidateCode(ctx context.Context, code, otpID string) (*models.OTP, error)
}

type otpService struct {
	otpRepo repositories.OTPRepository
	config  configs.Environment
}

func NewOTPService(otpRepo repositories.OTPRepository, config configs.Environment) OTPService {
	return &otpService{
		otpRepo: otpRepo,
		config:  config,
	}
}

func (s *otpService) CreateOTP(ctx context.Context, userID string) (*models.OTP, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("convert userID to ObjectID: %w", err)
	}

	otp := &models.OTP{
		UserID:    userIDObj,
		Code:      s.generateOTP(),
		ExpiresAt: time.Now().UTC().Add(s.config.OTP.ExpirationMinutes),
	}

	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("create otp: %w", err)
	}

	return otp, nil
}

func (s *otpService) ValidateCode(ctx context.Context, code, otpID string) (*models.OTP, error) {
	otp, err := s.otpRepo.FindByID(ctx, otpID)
	if err != nil {
		return nil, err
	}

	if otp.Code != code {
		return nil, models.ErrInvalidCode
	}

	if time.Now().After(otp.ExpiresAt) {
		return nil, models.ErrOTPExpired
	}

	if err := s.otpRepo.Delete(ctx, otpID); err != nil {
		return nil, fmt.Errorf("delete otp: %w", err)
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
