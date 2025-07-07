package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/scopes"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type ClientService interface {
	CreateClient(ctx context.Context, name string, description string, redirectURIs []string, grantTypes []string) (*models.ClientResponse, error)
	GetClientByClientID(ctx context.Context, clientID string) (*entities.Client, error)
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
	if err != nil && !errors.Is(err, domain.ErrClientNotFound) {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	if clientFromClientID != nil {
		return nil, fmt.Errorf("create client: %w", domain.ErrClientAlreadyExists)
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

func (s *clientService) GetClientByClientID(ctx context.Context, clientID string) (*entities.Client, error) {
	client, err := s.clientRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get client by client_id: %w", err)
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
