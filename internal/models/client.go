package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateClientPayload representa o payload para criação de cliente
type CreateClientPayload struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1,dive,uri"`
	GrantTypes   []string `json:"grant_types" validate:"required,min=1,dive,oneof=authorization_code refresh_token"`
}

// UpdateClientPayload representa o payload para atualização de cliente
type UpdateClientPayload struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=1"`
	Description  *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	RedirectURIs []string `json:"redirect_uris,omitempty" validate:"omitempty,min=1,dive,uri"`
	GrantTypes   []string `json:"grant_types,omitempty" validate:"omitempty,min=1,dive,oneof=authorization_code refresh_token"`
	Scopes       []string `json:"scopes,omitempty" validate:"omitempty,min=1"`
}

// ClientResponse representa a resposta da API para cliente
type ClientResponse struct {
	ID           primitive.ObjectID `json:"id"`
	ClientID     string             `json:"client_id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	RedirectURIs []string           `json:"redirect_uris"`
	Scopes       []string           `json:"scopes"`
	CreatedAt    time.Time          `json:"created_at"`
}

// ClientListResponse representa a resposta da API para listagem de clientes
type ClientListResponse struct {
	Clients []ClientResponse `json:"clients"`
	Total   int64            `json:"total"`
}
