package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
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

type CreateClientPayload struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	RedirectURIs []string `json:"redirect_uris" validate:"required"`
	Scopes       []string `json:"scopes" validate:"required"`
}

type ClientResponse struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	ClientID     string             `json:"client_id" bson:"client_id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	RedirectURIs []string           `bson:"redirect_uris" json:"redirect_uris"`
	Scopes       []string           `bson:"scopes" json:"scopes"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
