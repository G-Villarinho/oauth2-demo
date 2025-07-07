package models

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
