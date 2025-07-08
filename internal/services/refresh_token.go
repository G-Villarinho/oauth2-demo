package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type RefreshTokenService interface {
	CreateRefreshToken(ctx context.Context, userID, clientID string, scopes []string) (*entities.RefreshToken, error)
}

type refreshTokenService struct {
	refreshTokenRepo repositories.RefreshTokenRepository
	config           *configs.Environment
}

func NewRefreshTokenService(
	refreshTokenRepo repositories.RefreshTokenRepository,
	config *configs.Environment,
) RefreshTokenService {
	return &refreshTokenService{
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
	}
}

func (s *refreshTokenService) CreateRefreshToken(ctx context.Context, userID, clientID string, scopes []string) (*entities.RefreshToken, error) {
	tokenVerifier, err := generateSecureRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("generate secure random string: %w", err)
	}

	tokenHash := hashToken(tokenVerifier)

	refreshToken := &entities.RefreshToken{
		TokenHash: tokenHash,
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		ExpiresAt: time.Now().Add(s.config.Security.RefreshTokenExpirationHours),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	return refreshToken, nil
}

func hashToken(token string) string {
	hasher := sha256.New()

	hasher.Write([]byte(token))

	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
