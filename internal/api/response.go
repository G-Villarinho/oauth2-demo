package api

type APIErrorResponse struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Errors  []ValidationErrorDetail `json:"errors"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
}
