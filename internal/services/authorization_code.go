package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

const (
	authorizationCodeLength = 32
	authorizationCodeExpiry = 10 * time.Minute
)

type AuthorizationCodeService interface {
	CreateAuthorizationCode(ctx context.Context, input models.CreateAuthorizationCodeInput) (*entities.AuthorizationCode, error)
	ValidateAuthorizationCode(ctx context.Context, code string, codeVerifier string) (*entities.AuthorizationCode, error)
}

type authorizationCodeService struct {
	authorizationCodeRepo repositories.AuthorizationCodeRepository
}

func NewAuthorizationCodeService(
	authorizationCodeRepo repositories.AuthorizationCodeRepository) AuthorizationCodeService {
	return &authorizationCodeService{
		authorizationCodeRepo: authorizationCodeRepo,
	}
}

func (s *authorizationCodeService) CreateAuthorizationCode(ctx context.Context, input models.CreateAuthorizationCodeInput) (*entities.AuthorizationCode, error) {
	code, err := generateSecureRandomString(authorizationCodeLength)
	if err != nil {
		return nil, fmt.Errorf("generate secure random string: %w", err)
	}

	expiresAt := time.Now().Add(authorizationCodeExpiry)

	authorizationCode := &entities.AuthorizationCode{
		Code:                code,
		UserID:              input.UserID,
		ClientID:            input.ClientID,
		RedirectURI:         input.RedirectURI,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		ExpiresAt:           expiresAt,
		Scopes:              input.Scopes,
	}

	if err := s.authorizationCodeRepo.Create(ctx, authorizationCode); err != nil {
		return nil, fmt.Errorf("create authorization code: %w", err)
	}

	return authorizationCode, nil
}

func (s *authorizationCodeService) ValidateAuthorizationCode(ctx context.Context, code string, codeVerifier string) (*entities.AuthorizationCode, error) {
	authorizationCode, err := s.authorizationCodeRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("find authorization code: %w", err)
	}

	if authorizationCode.IsExpired() {
		return nil, domain.ErrAuthorizationCodeExpired
	}

	if !validatePKCE(codeVerifier, authorizationCode.CodeChallenge, authorizationCode.CodeChallengeMethod) {
		return nil, domain.ErrAuthorizationCodeInvalid
	}

	return authorizationCode, nil
}

func generateSecureRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func validatePKCE(verifier, challenge, method string) bool {
	switch method {
	case "S256":
		verifierHash := sha256.Sum256([]byte(verifier))
		calculatedChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(verifierHash[:])
		return calculatedChallenge == challenge
	case "plain":
		return verifier == challenge
	default:
		return false
	}
}
