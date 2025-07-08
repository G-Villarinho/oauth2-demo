package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id"`
	TokenHash string             `bson:"token_hash"`
	UserID    string             `bson:"user_id"`
	ClientID  string             `bson:"client_id"`
	Scopes    []string           `bson:"scopes"`
	ExpiresAt time.Time          `bson:"expires_at"`
	RevokedAt time.Time          `bson:"revoked_at"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (r *RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now().UTC())
}

func (r *RefreshToken) IsRevoked() bool {
	return r.RevokedAt.Before(time.Now().UTC())
}

func (r *RefreshToken) IsValidClientID(clientID string) bool {
	return r.ClientID == clientID
}
