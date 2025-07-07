package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUserPayload representa o payload para criação de usuário
type CreateUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// UpdateUserPayload representa o payload para atualização de usuário
type UpdateUserPayload struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=1"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse representa a resposta da API para usuário
type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty"`
}
