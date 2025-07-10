package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateAuthorizationCode(t *testing.T) {
	t.Run("should return success when valid input is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "user123"
		clientID := "client456"
		redirectURI := "https://example.com/callback"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		codeChallengeMethod := "S256"
		scopes := []string{"read", "write"}

		input := models.CreateAuthorizationCodeInput{
			UserID:              userID,
			ClientID:            clientID,
			RedirectURI:         redirectURI,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
			Scopes:              scopes,
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.AuthorizationCode")).
			Return(nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.CreateAuthorizationCode(ctx, input)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Code)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, clientID, result.ClientID)
		assert.Equal(t, redirectURI, result.RedirectURI)
		assert.Equal(t, codeChallenge, result.CodeChallenge)
		assert.Equal(t, codeChallengeMethod, result.CodeChallengeMethod)
		assert.Equal(t, scopes, result.Scopes)
		assert.True(t, result.ExpiresAt.After(time.Now()))
		assert.True(t, result.ExpiresAt.Before(time.Now().Add(11*time.Minute)))
	})

	t.Run("should return error when repository fails to create", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedError := errors.New("database connection failed")

		input := models.CreateAuthorizationCodeInput{
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
			CodeChallengeMethod: "S256",
			Scopes:              []string{"read"},
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.AuthorizationCode")).
			Return(expectedError)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.CreateAuthorizationCode(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create authorization code")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should generate unique codes for different calls", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := models.CreateAuthorizationCodeInput{
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
			CodeChallengeMethod: "S256",
			Scopes:              []string{"read"},
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.AuthorizationCode")).
			Return(nil).
			Times(2)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result1, err1 := service.CreateAuthorizationCode(ctx, input)
		result2, err2 := service.CreateAuthorizationCode(ctx, input)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotNil(t, result1)
		assert.NotNil(t, result2)
		assert.NotEqual(t, result1.Code, result2.Code)
	})

	t.Run("should handle empty scopes", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := models.CreateAuthorizationCodeInput{
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
			CodeChallengeMethod: "S256",
			Scopes:              []string{},
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.AuthorizationCode")).
			Return(nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.CreateAuthorizationCode(ctx, input)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Scopes)
	})
}

func TestValidateAuthorizationCode(t *testing.T) {
	t.Run("should return success when valid code and verifier are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code_123"
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		expiresAt := time.Now().Add(5 * time.Minute)

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "S256",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read", "write"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, authorizationCode, result)
	})

	t.Run("should return success when using plain PKCE method", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code_123"
		codeVerifier := "plain_verifier"
		codeChallenge := "plain_verifier"
		expiresAt := time.Now().Add(5 * time.Minute)

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "plain",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, authorizationCode, result)
	})

	t.Run("should return error when authorization code is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "invalid_code"
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(nil, domain.ErrAuthorizationCodeNotFound)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find authorization code")
		assert.Contains(t, err.Error(), domain.ErrAuthorizationCodeNotFound.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code"
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		expectedError := errors.New("database connection failed")

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(nil, expectedError)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find authorization code")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when authorization code is expired", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "expired_code"
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		expiresAt := time.Now().Add(-5 * time.Minute) // Expired 5 minutes ago

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "S256",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrAuthorizationCodeExpired, err)
	})

	t.Run("should return error when PKCE validation fails for S256 method", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code"
		codeVerifier := "invalid_verifier"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		expiresAt := time.Now().Add(5 * time.Minute)

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "S256",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrAuthorizationCodeInvalid, err)
	})

	t.Run("should return error when PKCE validation fails for plain method", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code"
		codeVerifier := "wrong_verifier"
		codeChallenge := "correct_challenge"
		expiresAt := time.Now().Add(5 * time.Minute)

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "plain",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrAuthorizationCodeInvalid, err)
	})

	t.Run("should return error when PKCE method is unsupported", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "valid_code"
		codeVerifier := "any_verifier"
		codeChallenge := "any_challenge"
		expiresAt := time.Now().Add(5 * time.Minute)

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                code,
			UserID:              "user123",
			ClientID:            "client456",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: "unsupported_method",
			ExpiresAt:           expiresAt,
			Scopes:              []string{"read"},
			CreatedAt:           time.Now(),
		}

		mockRepo := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockRepo.EXPECT().
			FindByCode(ctx, code).
			Return(authorizationCode, nil)

		service := NewAuthorizationCodeService(mockRepo)

		// Act
		result, err := service.ValidateAuthorizationCode(ctx, code, codeVerifier)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrAuthorizationCodeInvalid, err)
	})
}

// Testes auxiliares para as funções privadas
func TestGenerateSecureRandomString(t *testing.T) {
	t.Run("should generate strings of correct length", func(t *testing.T) {
		// Act
		result, err := generateSecureRandomString(32)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Len(t, result, 44) // base64.URLEncoding produces 44 chars for 32 bytes
	})

	t.Run("should generate different strings on each call", func(t *testing.T) {
		// Act
		result1, err1 := generateSecureRandomString(16)
		result2, err2 := generateSecureRandomString(16)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, result1, result2)
	})
}

func TestValidatePKCE(t *testing.T) {
	t.Run("should validate S256 method correctly", func(t *testing.T) {
		// Arrange
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		method := "S256"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.True(t, result)
	})

	t.Run("should validate plain method correctly", func(t *testing.T) {
		// Arrange
		codeVerifier := "plain_verifier"
		codeChallenge := "plain_verifier"
		method := "plain"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.True(t, result)
	})

	t.Run("should return false for invalid S256 challenge", func(t *testing.T) {
		// Arrange
		codeVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		codeChallenge := "invalid_challenge"
		method := "S256"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.False(t, result)
	})

	t.Run("should return false for invalid plain challenge", func(t *testing.T) {
		// Arrange
		codeVerifier := "correct_verifier"
		codeChallenge := "wrong_challenge"
		method := "plain"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.False(t, result)
	})

	t.Run("should return false for unsupported method", func(t *testing.T) {
		// Arrange
		codeVerifier := "any_verifier"
		codeChallenge := "any_challenge"
		method := "unsupported"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.False(t, result)
	})

	t.Run("should handle empty strings", func(t *testing.T) {
		// Arrange
		codeVerifier := ""
		codeChallenge := ""
		method := "plain"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.True(t, result)
	})

	t.Run("should validate S256 with custom verifier", func(t *testing.T) {
		// Arrange
		codeVerifier := "test_verifier_123"
		verifierHash := sha256.Sum256([]byte(codeVerifier))
		codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(verifierHash[:])
		method := "S256"

		// Act
		result := validatePKCE(codeVerifier, codeChallenge, method)

		// Assert
		assert.True(t, result)
	})
}
