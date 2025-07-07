package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/pkg/ecdsa"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateOTPTokenJWT(ctx context.Context, jti string, expiresAt time.Time) (string, error)
	GenerateAccessTokenJWT(ctx context.Context, firstName, lastName, email, userID string, expiresAt time.Time) (string, error)
}

type jwtService struct {
	ecdsa ecdsa.EcdsaKeyPair
}

func NewJWTService(ecdsa ecdsa.EcdsaKeyPair) JWTService {
	return &jwtService{
		ecdsa: ecdsa,
	}
}

func (s *jwtService) GenerateOTPTokenJWT(ctx context.Context, jti string, expiresAt time.Time) (string, error) {
	privateKey, err := s.ecdsa.ParseECDSAPrivateKey()
	if err != nil {
		return "", fmt.Errorf("parse ecdsa private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, models.OTPTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "aetheris-id",
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Audience:  jwt.ClaimStrings{"aetheris-id"},
		},
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) GenerateAccessTokenJWT(ctx context.Context, firstName, lastName, email, userID string, expiresAt time.Time) (string, error) {
	privateKey, err := s.ecdsa.ParseECDSAPrivateKey()
	if err != nil {
		return "", fmt.Errorf("parse ecdsa private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, models.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "aetheris-id",
			ID:        userID,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Audience:  jwt.ClaimStrings{"aetheris-id"},
		},
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}
