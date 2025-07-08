package models

import "github.com/golang-jwt/jwt/v5"

type OTPTokenClaims struct {
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"typ"`
}

type IDTokenClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"typ"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}
