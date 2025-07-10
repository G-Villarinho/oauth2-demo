package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateClient(t *testing.T) {
	t.Run("should return success when valid input is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client for testing purposes"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "test-client@aetheris-lab-connect").
			Return(nil, domain.ErrClientNotFound)

		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.Client")).
			Return(nil)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test-client@aetheris-lab-connect", result.ClientID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, description, result.Description)
		assert.Equal(t, redirectURIs, result.RedirectURIs)
		assert.NotEmpty(t, result.Scopes)
	})

	t.Run("should return error when client already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Existing Client"
		description := "A client that already exists"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}

		existingClient := &entities.Client{
			ID:           primitive.NewObjectID(),
			Name:         name,
			ClientID:     "existing-client@aetheris-lab-connect",
			Description:  description,
			RedirectURIs: redirectURIs,
			GrantTypes:   grantTypes,
		}

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "existing-client@aetheris-lab-connect").
			Return(existingClient, nil)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create client")
		assert.Contains(t, err.Error(), domain.ErrClientAlreadyExists.Error())
	})

	t.Run("should return error when repository fails to check existing client", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}
		expectedError := errors.New("database connection failed")

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "test-client@aetheris-lab-connect").
			Return(nil, expectedError)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when repository fails to create client", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{"authorization_code"}
		expectedError := errors.New("failed to create client")

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "test-client@aetheris-lab-connect").
			Return(nil, domain.ErrClientNotFound)

		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.Client")).
			Return(expectedError)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create client")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should handle empty redirect URIs", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{}
		grantTypes := []string{"authorization_code"}

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "test-client@aetheris-lab-connect").
			Return(nil, domain.ErrClientNotFound)

		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.Client")).
			Return(nil)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.RedirectURIs)
	})

	t.Run("should handle empty grant types", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "Test Client"
		description := "A test client"
		redirectURIs := []string{"https://example.com/callback"}
		grantTypes := []string{}

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, "test-client@aetheris-lab-connect").
			Return(nil, domain.ErrClientNotFound)

		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*entities.Client")).
			Return(nil)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.CreateClient(ctx, name, description, redirectURIs, grantTypes)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestGetClientByClientID(t *testing.T) {
	t.Run("should return client when valid client ID is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client@aetheris-lab-connect"

		expectedClient := &entities.Client{
			ID:           primitive.NewObjectID(),
			Name:         "Test Client",
			ClientID:     clientID,
			Description:  "A test client",
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code"},
		}

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.GetClientByClientID(ctx, clientID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedClient, result)
	})

	t.Run("should return error when client is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "nonexistent-client@aetheris-lab-connect"

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, clientID).
			Return(nil, domain.ErrClientNotFound)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.GetClientByClientID(ctx, clientID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
		assert.Contains(t, err.Error(), domain.ErrClientNotFound.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client@aetheris-lab-connect"
		expectedError := errors.New("database connection failed")

		mockRepo := mocks.NewClientRepositoryMock(t)
		mockRepo.EXPECT().
			GetByClientID(ctx, clientID).
			Return(nil, expectedError)

		service := NewClientService(mockRepo)

		// Act
		result, err := service.GetClientByClientID(ctx, clientID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get client by client_id")
		assert.Contains(t, err.Error(), expectedError.Error())
	})
}

func TestGenerateClientID(t *testing.T) {
	t.Run("should generate valid client ID from simple name", func(t *testing.T) {
		// Arrange
		name := "Test Client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "test-client@aetheris-lab-connect", result)
	})

	t.Run("should handle names with multiple spaces", func(t *testing.T) {
		// Arrange
		name := "My  Test   Client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "my--test---client@aetheris-lab-connect", result)
	})

	t.Run("should handle names with special characters", func(t *testing.T) {
		// Arrange
		name := "Test@Client#123"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient123@aetheris-lab-connect", result)
	})

	t.Run("should handle names with uppercase letters", func(t *testing.T) {
		// Arrange
		name := "UPPERCASE CLIENT"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "uppercase-client@aetheris-lab-connect", result)
	})

	t.Run("should handle names with numbers", func(t *testing.T) {
		// Arrange
		name := "Client 123"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "client-123@aetheris-lab-connect", result)
	})

	t.Run("should handle names with hyphens", func(t *testing.T) {
		// Arrange
		name := "test-client-name"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "test-client-name@aetheris-lab-connect", result)
	})

	t.Run("should handle names with mixed valid and invalid characters", func(t *testing.T) {
		// Arrange
		name := "Test!@#$%^&*()Client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient@aetheris-lab-connect", result)
	})

	t.Run("should handle empty name", func(t *testing.T) {
		// Arrange
		name := ""
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "@aetheris-lab-connect", result)
	})

	t.Run("should handle name with only special characters", func(t *testing.T) {
		// Arrange
		name := "!@#$%^&*()"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "@aetheris-lab-connect", result)
	})

	t.Run("should handle name with leading and trailing spaces", func(t *testing.T) {
		// Arrange
		name := "  Test Client  "
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "--test-client--@aetheris-lab-connect", result)
	})

	t.Run("should handle name with only numbers", func(t *testing.T) {
		// Arrange
		name := "123456"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "123456@aetheris-lab-connect", result)
	})

	t.Run("should handle name with only letters", func(t *testing.T) {
		// Arrange
		name := "TestClient"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient@aetheris-lab-connect", result)
	})

	t.Run("should handle name with consecutive hyphens", func(t *testing.T) {
		// Arrange
		name := "test--client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "test--client@aetheris-lab-connect", result)
	})

	t.Run("should handle name with unicode characters", func(t *testing.T) {
		// Arrange
		name := "TestáClient"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient@aetheris-lab-connect", result)
	})

	t.Run("should handle name with underscores", func(t *testing.T) {
		// Arrange
		name := "test_client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient@aetheris-lab-connect", result)
	})

	t.Run("should handle name with dots", func(t *testing.T) {
		// Arrange
		name := "test.client"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "testclient@aetheris-lab-connect", result)
	})

	t.Run("should handle name with mixed case and special characters", func(t *testing.T) {
		// Arrange
		name := "My-Test@Client#123"
		service := &clientService{}

		// Act
		result := service.generateClientID(name)

		// Assert
		assert.Equal(t, "my-testclient123@aetheris-lab-connect", result)
	})
}
