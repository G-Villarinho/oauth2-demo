package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/pkg/ecdsa"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService interface {
	GenerateOTPTokenJWT(ctx context.Context, jti string, expiresAt time.Time) (string, error)
	GenerateAccessTokenJWT(ctx context.Context, userID string, expiresAt time.Time) (string, error)
	GenerateIDTokenJWT(ctx context.Context, userID, name, email string, expiresAt time.Duration) (string, error)
	ValidateOTPTokenJWT(ctx context.Context, token string) (models.OTPTokenClaims, error)
	ValidateAccessTokenJWT(ctx context.Context, token string) (models.AccessTokenClaims, error)
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

func (s *jwtService) GenerateAccessTokenJWT(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	privateKey, err := s.ecdsa.ParseECDSAPrivateKey()
	if err != nil {
		return "", fmt.Errorf("parse ecdsa private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, models.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://id.aetheris-lab.com",
			ID:        primitive.NewObjectID().Hex(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Audience:  jwt.ClaimStrings{"https://app.aetheris-lab.com"},
			Subject:   userID,
		},
		TokenType: "Bearer",
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) GenerateIDTokenJWT(ctx context.Context, userID, name, email string, expiresAt time.Duration) (string, error) {
	privateKey, err := s.ecdsa.ParseECDSAPrivateKey()
	if err != nil {
		return "", fmt.Errorf("parse ecdsa private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, models.IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://id.aetheris-lab.com",
			ID:        primitive.NewObjectID().Hex(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			Audience:  jwt.ClaimStrings{"https://app.aetheris-lab.com"},
			Subject:   userID,
		},
		TokenType: "Bearer",
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) ValidateOTPTokenJWT(ctx context.Context, token string) (models.OTPTokenClaims, error) {
	publicKey, err := s.ecdsa.ParseECDSAPublicKey()
	if err != nil {
		return models.OTPTokenClaims{}, fmt.Errorf("parse ecdsa public key: %w", err)
	}

	claims := models.OTPTokenClaims{}
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return models.OTPTokenClaims{}, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}

func (s *jwtService) ValidateAccessTokenJWT(ctx context.Context, token string) (models.AccessTokenClaims, error) {
	publicKey, err := s.ecdsa.ParseECDSAPublicKey()
	if err != nil {
		return models.AccessTokenClaims{}, fmt.Errorf("parse ecdsa public key: %w", err)
	}

	claims := models.AccessTokenClaims{}
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return models.AccessTokenClaims{}, fmt.Errorf("parse token: %w", err)
	}

	return claims, nil
}
