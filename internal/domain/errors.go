package domain

import (
	"errors"
	"fmt"
	"time"
)

var (

	// OTP
	ErrOTPNotFound = errors.New("otp not found")
	ErrInvalidCode = errors.New("invalid code")
	ErrOTPExpired  = errors.New("otp expired")

	// Client
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidRedirectURI  = errors.New("invalid redirect uri")
	ErrInvalidGrantType    = errors.New("invalid grant type")
	ErrInvalidScope        = errors.New("invalid scope")

	// User
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")

	// Authorization Code
	ErrAuthorizationCodeNotFound = errors.New("authorization code not found")
	ErrAuthorizationCodeExpired  = errors.New("authorization code expired")
	ErrAuthorizationCodeInvalid  = errors.New("authorization code invalid")

	// OAuth
	ErrInvalidResponseType     = errors.New("invalid response type")
	ErrUnauthorizedClient      = errors.New("unauthorized client")
	ErrUnauthorizedRedirectURI = errors.New("unauthorized redirect uri")
	ErrUserAlreadyRegistered   = errors.New("user already registered")

	// Refresh Token
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")

	// ObjectID
	ErrInvalidObjectID = errors.New("invalid object id")
)

type ErrOTPNotResendable struct {
	TimeRemaining time.Duration
}

func (e *ErrOTPNotResendable) Error() string {
	return fmt.Sprintf("otp not resendable, time remaining: %s", e.TimeRemaining)
}
