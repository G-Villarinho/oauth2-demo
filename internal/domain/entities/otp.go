package entities

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrOTPNotFound = errors.New("otp not found")
	ErrInvalidCode = errors.New("invalid code")
	ErrOTPExpired  = errors.New("otp expired")
)

type OTP struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Code      string             `json:"code" bson:"code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// IsExpired verifica se o OTP expirou
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid verifica se o OTP é válido (não expirado)
func (o *OTP) IsValid() bool {
	return !o.IsExpired()
}

// ValidateCode verifica se o código fornecido é válido
func (o *OTP) ValidateCode(code string) error {
	if o.IsExpired() {
		return ErrOTPExpired
	}

	if o.Code != code {
		return ErrInvalidCode
	}

	return nil
}

// GetTimeUntilExpiration retorna o tempo restante até a expiração
func (o *OTP) GetTimeUntilExpiration() time.Duration {
	return o.ExpiresAt.Sub(time.Now())
}
