package models

import "time"

type LoginPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthenticatePayload struct {
	Code string `json:"code" validate:"required"`
}

type SendVerificationCodeResponse struct {
	OTPToken  string
	ExpiresAt time.Time
}

type AuthenticateResponse struct {
	AccessToken string
	ExpiresAt   time.Time
}
