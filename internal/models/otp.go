package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateOTPPayload representa o payload para criação de OTP
type CreateOTPPayload struct {
	UserID    primitive.ObjectID `json:"user_id" validate:"required"`
	Code      string             `json:"code" validate:"required"`
	ExpiresAt time.Time          `json:"expires_at" validate:"required"`
}

// ValidateOTPPayload representa o payload para validação de OTP
type ValidateOTPPayload struct {
	Code string `json:"code" validate:"required"`
}

// OTPResponse representa a resposta da API para OTP
type OTPResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	ExpiresAt time.Time          `json:"expires_at"`
	CreatedAt time.Time          `json:"created_at"`
}
