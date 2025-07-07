package models

type AuthorizePayload struct {
	ClientID            string   `query:"client_id" validate:"required"`
	RedirectURI         string   `query:"redirect_uri" validate:"required"`
	ResponseType        string   `query:"response_type" validate:"required"`
	Scope               []string `query:"scope" validate:"required"`
	State               string   `query:"state" validate:"required"`
	CodeChallenge       string   `query:"code_challenge" validate:"required"`
	CodeChallengeMethod string   `query:"code_challenge_method" validate:"required"`
}

type AuthorizeResponse struct {
	RedirectURL string
}
