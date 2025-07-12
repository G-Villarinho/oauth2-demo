package services

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/scopes"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type OAuthService interface {
	Authorize(ctx context.Context, input models.AuthorizeInput) (*models.AuthorizeResponse, error)
	ExchangeCodeForToken(ctx context.Context, input models.ExchangeAuthorizationCodeInput) (*models.TokenResponse, error)
}

type oauthService struct {
	clientService       ClientService
	authCodeService     AuthorizationCodeService
	jwtService          JWTService
	userRepo            repositories.UserRepository
	refreshTokenService RefreshTokenService
	config              *configs.Environment
}

func NewOAuthService(
	clientService ClientService,
	authCodeService AuthorizationCodeService,
	jwtService JWTService,
	userRepo repositories.UserRepository,
	refreshTokenService RefreshTokenService,
	config *configs.Environment,
) OAuthService {
	return &oauthService{
		clientService:       clientService,
		authCodeService:     authCodeService,
		jwtService:          jwtService,
		userRepo:            userRepo,
		refreshTokenService: refreshTokenService,
		config:              config,
	}
}

func (s *oauthService) Authorize(ctx context.Context, input models.AuthorizeInput) (*models.AuthorizeResponse, error) {
	if input.UserID == "" {
		loginRedirectURL, err := s.loginRedirectURL(input)
		if err != nil {
			return nil, fmt.Errorf("login redirect url: %w", err)
		}

		return &models.AuthorizeResponse{
			RedirectURL: loginRedirectURL,
		}, nil
	}

	client, err := s.clientService.GetClientByClientID(ctx, input.ClientID)
	if err != nil {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	if err := s.validateOAuthParameters(client, input); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	fmt.Println("input.scope", input.Scope)

	authorizationCodeInput := models.CreateAuthorizationCodeInput{
		UserID:              input.UserID,
		ClientID:            input.ClientID,
		RedirectURI:         input.RedirectURI,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		Scopes:              input.Scope,
	}
	authorizationCode, err := s.authCodeService.CreateAuthorizationCode(ctx, authorizationCodeInput)
	if err != nil {
		return nil, fmt.Errorf("create authorization code: %w", err)
	}

	callbackURL, err := s.callbackURL(authorizationCode.Code, input.State, input.RedirectURI)
	if err != nil {
		return nil, fmt.Errorf("callback url: %w", err)
	}

	return &models.AuthorizeResponse{
		RedirectURL: callbackURL,
	}, nil
}

func (s *oauthService) ExchangeCodeForToken(ctx context.Context, input models.ExchangeAuthorizationCodeInput) (*models.TokenResponse, error) {
	authorizationCode, err := s.authCodeService.ValidateAuthorizationCode(ctx, input.Code, input.CodeVerifier)
	if err != nil {
		return nil, fmt.Errorf("validate authorization code: %w", err)
	}

	if !authorizationCode.IsValidClientID(input.ClientID) {
		return nil, fmt.Errorf("exchange code for token %w: %s", domain.ErrUnauthorizedClient, input.ClientID)
	}

	if !authorizationCode.IsValidRedirectURI(input.RedirectURI) {
		return nil, fmt.Errorf("exchange code for token %w", domain.ErrUnauthorizedRedirectURI)
	}

	client, err := s.clientService.GetClientByClientID(ctx, authorizationCode.ClientID)
	if err != nil {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	hasRefreshToken := client.IsValidGrantType("refresh_token")
	accessTokenExpiresAt := time.Now().Add(s.config.Security.AccessTokenExpirationHours)
	if hasRefreshToken {
		accessTokenExpiresAt = time.Now().Add(s.config.Security.RefreshTokenExpirationHours)
	}

	accessToken, err := s.jwtService.GenerateAccessTokenJWT(ctx, authorizationCode.UserID, accessTokenExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	var refreshTokenHash string
	if hasRefreshToken {
		refreshToken, err := s.refreshTokenService.CreateRefreshToken(ctx, authorizationCode.UserID, client.ClientID, authorizationCode.Scopes)
		if err != nil {
			return nil, fmt.Errorf("create refresh token: %w", err)
		}

		refreshTokenHash = refreshToken.TokenHash
	}

	hasOpenID := scopes.HasScope(authorizationCode.Scopes, "openid")
	var idToken string
	if hasOpenID {
		user, err := s.userRepo.FindByID(ctx, authorizationCode.UserID)
		if err != nil {
			return nil, fmt.Errorf("find user by id: %w", err)
		}

		idToken, err = s.jwtService.GenerateIDTokenJWT(ctx, authorizationCode.UserID, user.GetFullName(), user.Email, s.config.Security.IDTokenExpirationMinutes)
		if err != nil {
			return nil, fmt.Errorf("generate id token: %w", err)
		}
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		IDToken:      idToken,
		RefreshToken: refreshTokenHash,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(accessTokenExpiresAt).Seconds()),
	}, nil
}

func (s *oauthService) loginRedirectURL(input models.AuthorizeInput) (string, error) {
	originalURL, err := s.authorizeURL(input)
	if err != nil {
		return "", fmt.Errorf("authorize url: %w", err)
	}

	loginURL, err := url.Parse(s.config.URLs.ClientLoginURL)
	if err != nil {
		return "", fmt.Errorf("invalid login url: %w", err)
	}

	query := loginURL.Query()
	query.Set("continue", originalURL)
	loginURL.RawQuery = query.Encode()

	return loginURL.String(), nil
}

func (s *oauthService) authorizeURL(input models.AuthorizeInput) (string, error) {
	baseURL, err := url.Parse(s.config.URLs.APIBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}

	baseURL.Path = "api/v1/oauth/authorize"

	query := baseURL.Query()
	query.Set("client_id", input.ClientID)
	query.Set("redirect_uri", input.RedirectURI)
	query.Set("response_type", input.ResponseType)
	query.Set("code_challenge", input.CodeChallenge)
	query.Set("code_challenge_method", input.CodeChallengeMethod)
	query.Set("scope", strings.Join(input.Scope, " "))
	query.Set("state", input.State)

	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil
}

func (s *oauthService) callbackURL(code, state, redirectURI string) (string, error) {
	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		return "", fmt.Errorf("invalid redirect uri: %w", err)
	}

	query := redirectURL.Query()
	query.Set("code", code)
	query.Set("state", state)

	redirectURL.RawQuery = query.Encode()

	return redirectURL.String(), nil
}

func (s *oauthService) validateOAuthParameters(client *entities.Client, input models.AuthorizeInput) error {
	if !client.IsValidRedirectURI(input.RedirectURI) {
		return fmt.Errorf("authorize: %w", domain.ErrInvalidRedirectURI)
	}

	if err := client.ValidateResponseType(input.ResponseType); err != nil {
		return err
	}

	if err := client.ValidateScopes(input.Scope); err != nil {
		return err
	}

	return nil
}
