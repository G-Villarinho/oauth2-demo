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
