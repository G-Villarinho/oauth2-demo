package entities

import (
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email     string             `json:"email" bson:"email"`
	Code      string             `json:"code" bson:"code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	ResendAt  *time.Time         `json:"resend_at" bson:"resend_at"`
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func (o *OTP) IsValid() bool {
	return !o.IsExpired()
}

func (o *OTP) IsResendable() bool {
	if o.ResendAt == nil {
		return true
	}

	return time.Since(*o.ResendAt) >= 60*time.Second
}

func (o *OTP) ValidateCode(code string) error {
	if o.IsExpired() {
		return domain.ErrOTPExpired
	}

	if o.Code != code {
		return domain.ErrInvalidCode
	}

	return nil
}

func (o *OTP) GetTimeUntilExpiration() time.Duration {
	return time.Until(o.ExpiresAt)
}
