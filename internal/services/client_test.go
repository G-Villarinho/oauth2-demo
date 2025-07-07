package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/scopes"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateClient(t *testing.T) {
	t.Run("should return success when valid client data is provided and client does not exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client for OAuth"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code", "refresh_token"}

		expectedClientID := "test-client@aetheris-lab-connect"
		expectedScopes := scopes.GetDefaultFirstPartyScopes()

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return client.Name == name &&
					client.Description == description &&
					client.ClientID == expectedClientID &&
					len(client.RedirectURIs) == 1 &&
					client.RedirectURIs[0] == redirectURIs[0] &&
					len(client.GrantTypes) == 2 &&
					client.GrantTypes[0] == grantTypes[0] &&
					client.GrantTypes[1] == grantTypes[1] &&
					len(client.Scopes) == len(expectedScopes)
			})).
			Run(func(ctx context.Context, client *models.Client) {
				// Set CreatedAt to ensure it's not zero in the response
				client.CreatedAt = time.Now().UTC()
			}).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, description, result.Description)
		assert.Equal(t, expectedClientID, result.ClientID)
		assert.Equal(t, redirectURIs, result.RedirectURIs)
		assert.Equal(t, expectedScopes, result.Scopes)
		assert.False(t, result.CreatedAt.IsZero())
	})

	t.Run("should return error when client with same name already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Existing Client"
		description := "A client that already exists"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedClientID := "existing-client@aetheris-lab-connect"
		existingClient := &models.Client{
			ID:           primitive.NewObjectID(),
			ClientID:     expectedClientID,
			Name:         name,
			Description:  description,
			RedirectURIs: redirectURIs,
			GrantTypes:   grantTypes,
			Scopes:       scopes.GetDefaultFirstPartyScopes(),
			CreatedAt:    time.Now(),
		}

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(existingClient, nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create client")
		assert.Contains(t, err.Error(), "client already exists")
	})

	t.Run("should return error when repository fails to get client by client_id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedClientID := "test-client@aetheris-lab-connect"
		expectedError := errors.New("database connection failed")

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(nil, expectedError)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("should return error when repository fails to create client", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedClientID := "test-client@aetheris-lab-connect"
		expectedError := errors.New("failed to insert document")

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*models.Client")).
			Return(expectedError)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create client")
		assert.Contains(t, err.Error(), "failed to insert document")
	})

	t.Run("should generate correct client_id from name with special characters", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "My Test Client 123!"
		description := "A test client with special characters"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedClientID := "my-test-client-123@aetheris-lab-connect"

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return client.ClientID == expectedClientID
			})).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedClientID, result.ClientID)
	})

	t.Run("should generate correct client_id from name with uppercase letters", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "UPPERCASE CLIENT"
		description := "A test client with uppercase letters"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedClientID := "uppercase-client@aetheris-lab-connect"

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, expectedClientID).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return client.ClientID == expectedClientID
			})).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedClientID, result.ClientID)
	})

	t.Run("should assign default first party scopes to new client", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		expectedScopes := scopes.GetDefaultFirstPartyScopes()

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, mock.AnythingOfType("string")).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return len(client.Scopes) == len(expectedScopes) &&
					client.Scopes[0] == expectedScopes[0] &&
					client.Scopes[1] == expectedScopes[1]
			})).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedScopes, result.Scopes)
	})

	t.Run("should handle multiple redirect URIs correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client with multiple redirect URIs"
		redirectURIs := []string{
			"https://example.com/callback",
			"https://app.example.com/oauth/callback",
			"http://localhost:3000/callback",
		}
		grantTypes := []string{"authorization_code"}

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, mock.AnythingOfType("string")).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return len(client.RedirectURIs) == 3 &&
					client.RedirectURIs[0] == redirectURIs[0] &&
					client.RedirectURIs[1] == redirectURIs[1] &&
					client.RedirectURIs[2] == redirectURIs[2]
			})).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, redirectURIs, result.RedirectURIs)
	})

	t.Run("should handle multiple grant types correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client with multiple grant types"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code", "refresh_token"}

		mockClientRepo := mocks.NewClientRepositoryMock(t)
		mockClientRepo.EXPECT().
			GetByClientID(ctx, mock.AnythingOfType("string")).
			Return(nil, models.ErrClientNotFound)

		mockClientRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(client *models.Client) bool {
				return len(client.GrantTypes) == 2 &&
					client.GrantTypes[0] == grantTypes[0] &&
					client.GrantTypes[1] == grantTypes[1]
			})).
			Return(nil)

		clientService := NewClientService(mockClientRepo)

		// Act
		result, err := clientService.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		// Note: GrantTypes is not included in ClientResponse, only in Client model
		assert.Equal(t, name, result.Name)
		assert.Equal(t, description, result.Description)
	})
}
