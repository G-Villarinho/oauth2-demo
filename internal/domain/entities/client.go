package entities

import (
	"fmt"
	"slices"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	ClientID     string             `json:"client_id" bson:"client_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	GrantTypes   []string           `bson:"grant_types" json:"grant_types"`
	RedirectURIs []string           `bson:"redirect_uris" json:"redirect_uris"`
	Scopes       []string           `bson:"scopes" json:"scopes"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (c *Client) IsValidRedirectURI(redirectURI string) bool {
	return slices.Contains(c.RedirectURIs, redirectURI)
}

func (c *Client) IsValidGrantType(grantType string) bool {
	return slices.Contains(c.GrantTypes, grantType)
}

func (c *Client) IsValidScope(scope string) bool {
	return slices.Contains(c.Scopes, scope)
}

func (c *Client) HasScope(scope string) bool {
	return c.IsValidScope(scope)
}

func (c *Client) GetValidGrantTypes() []string {
	return c.GrantTypes
}

func (c *Client) GetValidRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) ValidateRedirectURI(redirectURI string) error {
	if !c.IsValidRedirectURI(redirectURI) {
		return fmt.Errorf("%w: %s", domain.ErrInvalidRedirectURI, redirectURI)
	}
	return nil
}

func (c *Client) ValidateResponseType(responseType string) error {
	switch responseType {
	case "code":
		if !c.IsValidGrantType("authorization_code") {
			return fmt.Errorf("%w: client does not support authorization_code grant type", domain.ErrInvalidGrantType)
		}
	case "token":
		if !c.IsValidGrantType("implicit") {
			return fmt.Errorf("%w: client does not support implicit grant type", domain.ErrInvalidGrantType)
		}
	default:
		return fmt.Errorf("%w: %s", domain.ErrInvalidResponseType, responseType)
	}
	return nil
}

func (c *Client) ValidateScopes(scopes []string) error {
	for _, scope := range scopes {
		if !c.IsValidScope(scope) {
			return fmt.Errorf("%w: %s", domain.ErrInvalidScope, scope)
		}
	}
	return nil
}
