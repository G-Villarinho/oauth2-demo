package entities

import (
	"errors"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidRedirectURI  = errors.New("invalid redirect uri")
	ErrInvalidGrantType    = errors.New("invalid grant type")
	ErrInvalidScope        = errors.New("invalid scope")
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

// IsValidRedirectURI verifica se a URI de redirecionamento é válida
func (c *Client) IsValidRedirectURI(redirectURI string) bool {
	return slices.Contains(c.RedirectURIs, redirectURI)
}

// IsValidGrantType verifica se o tipo de grant é válido
func (c *Client) IsValidGrantType(grantType string) bool {
	return slices.Contains(c.GrantTypes, grantType)
}

// IsValidScope verifica se o escopo é válido
func (c *Client) IsValidScope(scope string) bool {
	return slices.Contains(c.Scopes, scope)
}

// HasScope verifica se o cliente tem um escopo específico
func (c *Client) HasScope(scope string) bool {
	return c.IsValidScope(scope)
}

// GetValidGrantTypes retorna os tipos de grant válidos
func (c *Client) GetValidGrantTypes() []string {
	return c.GrantTypes
}

// GetValidRedirectURIs retorna as URIs de redirecionamento válidas
func (c *Client) GetValidRedirectURIs() []string {
	return c.RedirectURIs
}
