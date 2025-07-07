package models

import "github.com/golang-jwt/jwt/v5"

type OTPTokenClaims struct {
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	FirstName string
	LastName  string
	Email     string
}

func (a *AccessTokenClaims) GetFullName() string {
	return a.FirstName + " " + a.LastName
}
