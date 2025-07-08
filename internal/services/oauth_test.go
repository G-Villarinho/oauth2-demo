package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestExchangeCodeForToken(t *testing.T) {
	t.Run("should return token response when all validations pass", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{
			Security: configs.Security{
				AccessTokenExpirationHours:  1,
				RefreshTokenExpirationHours: 24,
				IDTokenExpirationMinutes:    15,
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

		input := models.ExchangeAuthorizationCodeInput{
			Code:         "valid-code",
			CodeVerifier: "valid-verifier",
			ClientID:     "test-client-id",
			RedirectURI:  "https://example.com/callback",
		}

		authCode := &entities.AuthorizationCode{
			UserID:      "test-user-id",
			ClientID:    "test-client-id",
			RedirectURI: "https://example.com/callback",
			Scopes:      []string{"openid", "profile", "email"},
		}

		client := &entities.Client{
			ClientID:   "test-client-id",
			GrantTypes: []string{"authorization_code", "refresh_token"},
		}

		user := &entities.User{
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
		}

		refreshToken := &entities.RefreshToken{
			TokenHash: "hashed-refresh-token",
		}

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, input.Code, input.CodeVerifier).Return(authCode, nil)
		mockClientService.EXPECT().GetClientByClientID(ctx, authCode.ClientID).Return(client, nil)
		mockJWTService.EXPECT().GenerateAccessTokenJWT(ctx, authCode.UserID, mock.AnythingOfType("time.Time")).Return("new-access-token", nil)
		mockRefreshTokenService.EXPECT().CreateRefreshToken(ctx, authCode.UserID, client.ClientID, authCode.Scopes).Return(refreshToken, nil)
		mockUserRepo.EXPECT().FindByID(ctx, authCode.UserID).Return(user, nil)
		mockJWTService.EXPECT().GenerateIDTokenJWT(ctx, authCode.UserID, user.GetFullName(), user.Email, config.Security.IDTokenExpirationMinutes).Return("new-id-token", nil)

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, input)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new-access-token", result.AccessToken)
		assert.Equal(t, "hashed-refresh-token", result.RefreshToken)
		assert.Equal(t, "new-id-token", result.IDToken)
		assert.Equal(t, "Bearer", result.TokenType)
		assert.InDelta(t, int(config.Security.RefreshTokenExpirationHours.Seconds()), result.ExpiresIn, 1)
	})

	t.Run("should return error when authorization code is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{}

		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		oauthService := NewOAuthService(nil, mockAuthCodeService, nil, nil, nil, config)

		input := models.ExchangeAuthorizationCodeInput{Code: "invalid-code"}
		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "invalid-code", "").Return(nil, errors.New("invalid code"))

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate authorization code")
	})

	t.Run("should return error when client ID is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		oauthService := NewOAuthService(nil, mockAuthCodeService, nil, nil, nil, config)

		input := models.ExchangeAuthorizationCodeInput{ClientID: "wrong-client-id", Code: "valid-code"}
		authCode := &entities.AuthorizationCode{ClientID: "correct-client-id"}

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "valid-code", "").Return(authCode, nil)

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUnauthorizedClient)
	})

	t.Run("should return error when redirect URI is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		oauthService := NewOAuthService(nil, mockAuthCodeService, nil, nil, nil, config)

		input := models.ExchangeAuthorizationCodeInput{RedirectURI: "wrong-uri", ClientID: "client-id", Code: "valid-code"}
		authCode := &entities.AuthorizationCode{ClientID: "client-id", RedirectURI: "correct-uri"}

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "valid-code", "").Return(authCode, nil)

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, input)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUnauthorizedRedirectURI)
	})

	t.Run("should return error when client service fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockClientService := mocks.NewClientServiceMock(t)
		oauthService := NewOAuthService(mockClientService, mockAuthCodeService, nil, nil, nil, config)

		authCode := &entities.AuthorizationCode{ClientID: "client-id", RedirectURI: "uri", Code: "code"}
		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "code", "").Return(authCode, nil)
		mockClientService.EXPECT().GetClientByClientID(ctx, "client-id").Return(nil, errors.New("db error"))

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, models.ExchangeAuthorizationCodeInput{Code: "code", ClientID: "client-id", RedirectURI: "uri"})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
	})

	t.Run("should return tokens without refresh token if client does not support it", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{Security: configs.Security{AccessTokenExpirationHours: 1}}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockClientService := mocks.NewClientServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		oauthService := NewOAuthService(mockClientService, mockAuthCodeService, mockJWTService, nil, nil, config)

		authCode := &entities.AuthorizationCode{UserID: "user-id", ClientID: "client-id", RedirectURI: "uri", Scopes: []string{}}
		client := &entities.Client{GrantTypes: []string{"authorization_code"}} // No refresh_token

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "code", "").Return(authCode, nil)
		mockClientService.EXPECT().GetClientByClientID(ctx, "client-id").Return(client, nil)
		mockJWTService.EXPECT().GenerateAccessTokenJWT(ctx, "user-id", mock.AnythingOfType("time.Time")).Return("access-token", nil)

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, models.ExchangeAuthorizationCodeInput{Code: "code", ClientID: "client-id", RedirectURI: "uri"})

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "access-token", result.AccessToken)
		assert.Empty(t, result.RefreshToken)
	})

	t.Run("should return error when generating refresh token fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{Security: configs.Security{RefreshTokenExpirationHours: 24}}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockClientService := mocks.NewClientServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockRefreshTokenService := mocks.NewRefreshTokenServiceMock(t)
		oauthService := NewOAuthService(mockClientService, mockAuthCodeService, mockJWTService, nil, mockRefreshTokenService, config)

		authCode := &entities.AuthorizationCode{UserID: "user-id", ClientID: "client-id", RedirectURI: "uri", Scopes: []string{}}
		client := &entities.Client{ClientID: "client-id", GrantTypes: []string{"refresh_token"}}

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "code", "").Return(authCode, nil)
		mockClientService.EXPECT().GetClientByClientID(ctx, "client-id").Return(client, nil)
		mockJWTService.EXPECT().GenerateAccessTokenJWT(ctx, "user-id", mock.AnythingOfType("time.Time")).Return("access-token", nil)
		mockRefreshTokenService.EXPECT().CreateRefreshToken(ctx, "user-id", "client-id", []string{}).Return(nil, errors.New("failed to create refresh token"))

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, models.ExchangeAuthorizationCodeInput{Code: "code", ClientID: "client-id", RedirectURI: "uri"})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create refresh token")
	})

	t.Run("should return error when user lookup for id_token fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		config := &configs.Environment{Security: configs.Security{AccessTokenExpirationHours: 1}}
		mockAuthCodeService := mocks.NewAuthorizationCodeServiceMock(t)
		mockClientService := mocks.NewClientServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		mockUserRepo := mocks.NewUserRepositoryMock(t)
		oauthService := NewOAuthService(mockClientService, mockAuthCodeService, mockJWTService, mockUserRepo, nil, config)

		authCode := &entities.AuthorizationCode{UserID: "user-id", ClientID: "client-id", RedirectURI: "uri", Scopes: []string{"openid"}}
		client := &entities.Client{GrantTypes: []string{}}

		mockAuthCodeService.EXPECT().ValidateAuthorizationCode(ctx, "code", "").Return(authCode, nil)
		mockClientService.EXPECT().GetClientByClientID(ctx, "client-id").Return(client, nil)
		mockJWTService.EXPECT().GenerateAccessTokenJWT(ctx, "user-id", mock.AnythingOfType("time.Time")).Return("access-token", nil)
		mockUserRepo.EXPECT().FindByID(ctx, "user-id").Return(nil, errors.New("user not found"))

		// Act
		result, err := oauthService.ExchangeCodeForToken(ctx, models.ExchangeAuthorizationCodeInput{Code: "code", ClientID: "client-id", RedirectURI: "uri"})

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by id")
	})
}
