package models

type CreateAuthorizationCodeInput struct {
	UserID              string
	ClientID            string
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod string
	Scopes              []string
}
