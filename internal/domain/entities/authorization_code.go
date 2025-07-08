package entities

import (
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthorizationCode struct {
	ID                  primitive.ObjectID `bson:"_id"`
	Code                string             `bson:"code"`
	UserID              string             `bson:"user_id"`
	ClientID            string             `bson:"client_id"`
	RedirectURI         string             `bson:"redirect_uri"`
	CodeChallenge       string             `bson:"code_challenge"`
	CodeChallengeMethod string             `bson:"code_challenge_method"`
	ExpiresAt           time.Time          `bson:"expires_at"`
	Scopes              []string           `bson:"scopes"`
	CreatedAt           time.Time          `bson:"created_at"`
}

func (a *AuthorizationCode) IsExpired() bool {
	return a.ExpiresAt.Before(time.Now().UTC())
}

func (a *AuthorizationCode) HasScope(scope string) bool {
	return slices.Contains(a.Scopes, scope)
}

func (a *AuthorizationCode) IsValidRedirectURI(redirectURI string) bool {
	return a.RedirectURI == redirectURI
}

func (a *AuthorizationCode) IsValidClientID(clientID string) bool {
	return a.ClientID == clientID
}
