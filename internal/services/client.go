package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
)

type ClientService interface {
	CreateClient(ctx context.Context, name string, description string, redirectURIs []string) (*models.ClientResponse, error)
	ValidateClient(ctx context.Context, clientID string, redirectURI string) (*models.Client, error)
}

type clientService struct {
	clientRepo repositories.ClientRepository
}

func NewClientService(clientRepo repositories.ClientRepository) ClientService {
	return &clientService{
		clientRepo: clientRepo,
	}
}

func (s *clientService) CreateClient(ctx context.Context, name string, description string, redirectURIs []string) (*models.ClientResponse, error) {
	clientId := s.generateClientID(name)

	clientFromClientID, err := s.clientRepo.GetByClientID(ctx, clientId)
	if err != nil {
		return nil, fmt.Errorf("get client by client_id: %w", err)
	}

	if clientFromClientID != nil {
		return nil, fmt.Errorf("create client: %w", models.ErrClientAlreadyExists)
	}

	client := &models.Client{
		Name:         name,
		Description:  description,
		RedirectURIs: redirectURIs,
		ClientID:     clientId,
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &models.ClientResponse{
		ID:           client.ID,
		ClientID:     client.ClientID,
		Name:         client.Name,
		Description:  client.Description,
		RedirectURIs: client.RedirectURIs,
		CreatedAt:    client.CreatedAt,
	}, nil
}

// ValidateClient implements ClientService.
func (s *clientService) ValidateClient(ctx context.Context, clientID string, redirectURI string) (*models.Client, error) {
	panic("unimplemented")
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
