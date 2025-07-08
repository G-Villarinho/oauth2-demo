package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthorize(t *testing.T) {
	t.Run("should return login redirect URL when user ID is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "", // Empty UserID
		}

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.RedirectURL)
		assert.Contains(t, result.RedirectURL, "https://example.com/login")
		assert.Contains(t, result.RedirectURL, "continue=")
	})

	t.Run("should return callback URL when user ID is provided and all validations pass", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"},
			Scopes:       []string{"read", "write"},
		}

		authorizationCode := &entities.AuthorizationCode{
			ID:                  primitive.NewObjectID(),
			Code:                "test-auth-code",
			UserID:              "test-user-id",
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scopes:              []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		mockAuthCodeService.EXPECT().
			CreateAuthorizationCode(ctx, models.CreateAuthorizationCodeInput{
				UserID:              "test-user-id",
				ClientID:            "test-client-id",
				RedirectURI:         "https://example.com/callback",
				CodeChallenge:       "test-challenge",
				CodeChallengeMethod: "S256",
				Scopes:              []string{"read", "write"},
			}).
			Return(authorizationCode, nil)

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.RedirectURL)
		assert.Contains(t, result.RedirectURL, "https://example.com/callback")
		assert.Contains(t, result.RedirectURL, "code=test-auth-code")
		assert.Contains(t, result.RedirectURL, "state=test-state")
	})

	t.Run("should return error when client is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "invalid-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "invalid-client-id").
			Return(nil, errors.New("client not found"))

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
	})

	t.Run("should return error when OAuth parameters validation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://invalid-redirect.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"}, // Different redirect URI
			Scopes:       []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate request")
	})

	t.Run("should return error when creating authorization code fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"},
			Scopes:       []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		mockAuthCodeService.EXPECT().
			CreateAuthorizationCode(ctx, models.CreateAuthorizationCodeInput{
				UserID:              "test-user-id",
				ClientID:            "test-client-id",
				RedirectURI:         "https://example.com/callback",
				CodeChallenge:       "test-challenge",
				CodeChallengeMethod: "S256",
				Scopes:              []string{"read", "write"},
			}).
			Return(nil, errors.New("failed to create authorization code"))

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create authorization code")
	})

	t.Run("should return error when redirect URI validation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://malicious.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"},
			Scopes:       []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate request")
	})

	t.Run("should return error when response type validation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "invalid_response_type",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"read", "write"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"},
			Scopes:       []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate request")
	})

	t.Run("should return error when scope validation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			URLs: configs.URLs{
				ClientLoginURL: "https://example.com/login",
				APIBaseURL:     "https://api.example.com",
			},
		}

		mockClientService := mocks.NewClientServiceMock(t)
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)

		oauthService := NewOAuthService(
			mockClientService,
			mockAuthCodeService,
			mockJWTService,
			mockUserRepo,
			mockRefreshTokenService,
			config,
		)

		input := models.AuthorizeInput{
			ClientID:            "test-client-id",
			RedirectURI:         "https://example.com/callback",
			ResponseType:        "code",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
			Scope:               []string{"invalid_scope"},
			State:               "test-state",
			UserID:              "test-user-id",
		}

		client := &entities.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     "test-client-id",
			Name:         "Test Client",
			GrantTypes:   []string{"authorization_code"},
			RedirectURIs: []string{"https://example.com/callback"},
			Scopes:       []string{"read", "write"},
		}

		mockClientService.EXPECT().
			GetClientByClientID(ctx, "test-client-id").
			Return(client, nil)

		// Act
		result, err := oauthService.Authorize(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate request")
	})
}
