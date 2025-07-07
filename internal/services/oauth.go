package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
)

type OAuthService interface {
	Authorize(ctx context.Context, input models.AuthorizeInput) (*models.AuthorizeResponse, error)
}

type oauthService struct {
	clientService            clientService
	authorizationCodeService authorizationCodeService
	jwtService               jwtService
	config                   *configs.Environment
}

func NewOAuthService(
	clientService clientService,
	authorizationCodeService authorizationCodeService,
	jwtService jwtService,
	config *configs.Environment,
) OAuthService {
	return &oauthService{
		clientService:            clientService,
		authorizationCodeService: authorizationCodeService,
		jwtService:               jwtService,
		config:                   config,
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

	if err := client.ValidateRequest(input.RedirectURI, input.ResponseType, input.Scope); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	return nil, nil
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

	baseURL.Path = "/authorize"

	query := baseURL.Query()
	query.Set("client_id", input.ClientID)
	query.Set("redirect_uri", input.RedirectURI)
	query.Set("response_type", input.ResponseType)
	query.Set("code_challenge", input.CodeChallenge)
	query.Set("code_challenge_method", input.CodeChallengeMethod)
	query.Set("state", input.State)

	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil
}
