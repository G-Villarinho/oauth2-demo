package domain

import "errors"

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
	ErrInvalidResponseType = errors.New("invalid response type")
)
