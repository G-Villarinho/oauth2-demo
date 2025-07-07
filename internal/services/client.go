package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/scopes"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type ClientService interface {
	CreateClient(ctx context.Context, name string, description string, redirectURIs []string, grantTypes []string) (*models.ClientResponse, error)
	ValidateClient(ctx context.Context, clientID string, redirectURI string, grantTypes []string, scopes []string) (*entities.Client, error)
}

type clientService struct {
	clientRepo repositories.ClientRepository
}

func NewClientService(clientRepo repositories.ClientRepository) ClientService {
	return &clientService{
		clientRepo: clientRepo,
	}
}

func (s *clientService) CreateClient(ctx context.Context, name string, description string, redirectURIs []string, grantTypes []string) (*models.ClientResponse, error) {
	clientId := s.generateClientID(name)

	clientFromClientID, err := s.clientRepo.GetByClientID(ctx, clientId)
	if err != nil && !errors.Is(err, entities.ErrClientNotFound) {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	if clientFromClientID != nil {
		return nil, fmt.Errorf("create client: %w", entities.ErrClientAlreadyExists)
	}

	client := &entities.Client{
		Name:         name,
		Description:  description,
		RedirectURIs: redirectURIs,
		ClientID:     clientId,
		Scopes:       scopes.GetDefaultFirstPartyScopes(),
		GrantTypes:   grantTypes,
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return models.ClientToResponse(client), nil
}

func (s *clientService) ValidateClient(ctx context.Context, clientID string, redirectURI string, grantTypes []string, scopes []string) (*entities.Client, error) {
	client, err := s.clientRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	if client == nil {
		return nil, fmt.Errorf("validate client: %w", entities.ErrClientNotFound)
	}

	if err := s.AllowsRequest(client, redirectURI, grantTypes, scopes); err != nil {
		return nil, fmt.Errorf("validate request parameters: %w", err)
	}

	return client, nil
}

func (s *clientService) generateClientID(name string) string {
	clientID := strings.ToLower(name)
	clientID = strings.ReplaceAll(clientID, " ", "-")

	var cleanID strings.Builder
	for _, char := range clientID {
		isValidChar := (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-'
		if isValidChar {
			cleanID.WriteRune(char)
		}
	}

	return cleanID.String() + "@aetheris-lab-connect"
}

func (s *clientService) AllowsRequest(client *entities.Client, redirectURI string, grantTypes []string, scopes []string) error {
	if !client.IsValidRedirectURI(redirectURI) {
		return fmt.Errorf("%w: %s", entities.ErrInvalidRedirectURI, redirectURI)
	}

	for _, grantType := range grantTypes {
		if !client.IsValidGrantType(grantType) {
			return fmt.Errorf("%w: %s", entities.ErrInvalidGrantType, grantType)
		}
	}

	for _, scope := range scopes {
		if !client.IsValidScope(scope) {
			return fmt.Errorf("%w: %s", entities.ErrInvalidScope, scope)
		}
	}

	return nil
}
