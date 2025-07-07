package models

import (
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserToResponse converte uma entidade User para UserResponse
func UserToResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// CreateUserPayloadToEntity converte CreateUserPayload para entidade User
func CreateUserPayloadToEntity(payload *CreateUserPayload) *entities.User {
	now := time.Now()
	return &entities.User{
		ID:        primitive.NewObjectID(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		CreatedAt: now,
		UpdatedAt: nil,
	}
}

// ClientToResponse converte uma entidade Client para ClientResponse
func ClientToResponse(client *entities.Client) *ClientResponse {
	return &ClientResponse{
		ID:           client.ID,
		ClientID:     client.ClientID,
		Name:         client.Name,
		Description:  client.Description,
		RedirectURIs: client.RedirectURIs,
		Scopes:       client.Scopes,
		CreatedAt:    client.CreatedAt,
	}
}

// CreateClientPayloadToEntity converte CreateClientPayload para entidade Client
func CreateClientPayloadToEntity(payload *CreateClientPayload) *entities.Client {
	now := time.Now()
	return &entities.Client{
		ID:           primitive.NewObjectID(),
		ClientID:     generateClientID(), // Função que você precisará implementar
		Name:         payload.Name,
		Description:  payload.Description,
		GrantTypes:   payload.GrantTypes,
		RedirectURIs: payload.RedirectURIs,
		Scopes:       []string{}, // Escopos vazios por padrão
		CreatedAt:    now,
		UpdatedAt:    nil,
	}
}

// OTPToResponse converte uma entidade OTP para OTPResponse
func OTPToResponse(otp *entities.OTP) *OTPResponse {
	return &OTPResponse{
		ID:        otp.ID,
		UserID:    otp.UserID,
		ExpiresAt: otp.ExpiresAt,
		CreatedAt: otp.CreatedAt,
	}
}

// CreateOTPPayloadToEntity converte CreateOTPPayload para entidade OTP
func CreateOTPPayloadToEntity(payload *CreateOTPPayload) *entities.OTP {
	now := time.Now()
	return &entities.OTP{
		ID:        primitive.NewObjectID(),
		UserID:    payload.UserID,
		Code:      payload.Code,
		ExpiresAt: payload.ExpiresAt,
		CreatedAt: now,
	}
}

// generateClientID gera um ID único para o cliente
// Esta é uma implementação simples - você pode melhorar conforme necessário
func generateClientID() string {
	return "client_" + primitive.NewObjectID().Hex()[:12]
}
