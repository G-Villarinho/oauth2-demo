package models

import "strings"

type AuthorizePayload struct {
	ClientID            string `query:"client_id" validate:"required"`
	RedirectURI         string `query:"redirect_uri" validate:"required"`
	ResponseType        string `query:"response_type" validate:"required"`
	Scope               string `query:"scope" validate:"required"`
	State               string `query:"state" validate:"required"`
	CodeChallenge       string `query:"code_challenge" validate:"required"`
	CodeChallengeMethod string `query:"code_challenge_method" validate:"required"`
}

type ExchangeAuthorizationCodePayload struct {
	Code         string `form:"code"`
	CodeVerifier string `form:"code_verifier"`
	ClientID     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
}

type AuthorizeInput struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	CodeChallenge       string
	Scope               []string
	CodeChallengeMethod string
	State               string
	UserID              string
}

type ExchangeAuthorizationCodeInput struct {
	Code         string
	CodeVerifier string
	ClientID     string
	RedirectURI  string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthorizeResponse struct {
	RedirectURL string
}

func NewAuthorizeInput(payload AuthorizePayload, userID string) AuthorizeInput {
	return AuthorizeInput{
		ClientID:            payload.ClientID,
		RedirectURI:         payload.RedirectURI,
		ResponseType:        payload.ResponseType,
		CodeChallenge:       payload.CodeChallenge,
		CodeChallengeMethod: payload.CodeChallengeMethod,
		State:               payload.State,
		Scope:               strings.Split(payload.Scope, " "),
		UserID:              userID,
	}
}

func NewExchangeAuthorizationCodeInput(payload ExchangeAuthorizationCodePayload) ExchangeAuthorizationCodeInput {
	return ExchangeAuthorizationCodeInput{
		Code:         payload.Code,
		CodeVerifier: payload.CodeVerifier,
		ClientID:     payload.ClientID,
		RedirectURI:  payload.RedirectURI,
	}
}
