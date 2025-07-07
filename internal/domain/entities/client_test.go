package entities

import (
	"testing"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ValidateRequest(t *testing.T) {
	t.Run("should return success when all parameters are valid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code", "refresh_token"},
			Scopes:       []string{"openid", "profile", "email"},
		}

		// Act
		err := client.ValidateRequest("https://example.com/callback", []string{"authorization_code"}, []string{"openid"})

		// Assert
		require.NoError(t, err)
	})

	t.Run("should return error when redirect URI is invalid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code"},
			Scopes:       []string{"openid"},
		}

		// Act
		err := client.ValidateRequest("https://malicious.com/callback", []string{"authorization_code"}, []string{"openid"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidRedirectURI)
		assert.Contains(t, err.Error(), "https://malicious.com/callback")
	})

	t.Run("should return error when grant type is invalid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code"},
			Scopes:       []string{"openid"},
		}

		// Act
		err := client.ValidateRequest("https://example.com/callback", []string{"invalid_grant_type"}, []string{"openid"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidGrantType)
		assert.Contains(t, err.Error(), "invalid_grant_type")
	})

	t.Run("should return error when scope is invalid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code"},
			Scopes:       []string{"openid"},
		}

		// Act
		err := client.ValidateRequest("https://example.com/callback", []string{"authorization_code"}, []string{"invalid_scope"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidScope)
		assert.Contains(t, err.Error(), "invalid_scope")
	})

	t.Run("should return error when multiple grant types are provided and one is invalid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code", "refresh_token"},
			Scopes:       []string{"openid"},
		}

		// Act
		err := client.ValidateRequest("https://example.com/callback", []string{"authorization_code", "invalid_grant"}, []string{"openid"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidGrantType)
		assert.Contains(t, err.Error(), "invalid_grant")
	})

	t.Run("should return error when multiple scopes are provided and one is invalid", func(t *testing.T) {
		// Arrange
		client := &Client{
			RedirectURIs: []string{"https://example.com/callback"},
			GrantTypes:   []string{"authorization_code"},
			Scopes:       []string{"openid", "profile"},
		}

		// Act
		err := client.ValidateRequest("https://example.com/callback", []string{"authorization_code"}, []string{"openid", "invalid_scope"})

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidScope)
		assert.Contains(t, err.Error(), "invalid_scope")
	})
}
