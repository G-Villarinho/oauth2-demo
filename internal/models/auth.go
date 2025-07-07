package models

type LoginPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthenticatePayload struct {
	Code string `json:"code" validate:"required"`
}

type SendVerificationCodeResponse struct {
	OTPToken string
}

type AuthenticateResponse struct {
	AccessToken string
}
