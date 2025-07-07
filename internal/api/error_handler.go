// internal/api/error_handler.go
package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, ectx echo.Context) {
	if ectx.Response().Committed {
		return
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		details := make([]ValidationErrorDetail, len(validationErrs))
		for i, fieldErr := range validationErrs {
			details[i] = ValidationErrorDetail{
				Field:   strings.ToLower(fieldErr.Field()),
				Message: formatValidationMessage(fieldErr),
				Tag:     fieldErr.Tag(),
				Value:   fmt.Sprintf("%v", fieldErr.Value()),
			}
		}

		response := APIErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "A validação dos dados falhou. Verifique os erros e tente novamente.",
			Errors:  details,
		}

		slog.Warn("Validation error", "details", details)

		if err := ectx.JSON(http.StatusUnprocessableEntity, response); err != nil {
			slog.Error("Failed to write validation error response", "error", err)
		}
		return
	}

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		response := APIErrorResponse{
			Code:    fmt.Sprintf("HTTP_%d", httpErr.Code),
			Message: fmt.Sprintf("%v", httpErr.Message),
		}

		if err := ectx.JSON(httpErr.Code, response); err != nil {
			slog.Error("Failed to write HTTP error response", "error", err)
		}
		return
	}

	slog.Error("Unhandled internal error", "error", err.Error())
	response := APIErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Ocorreu um erro interno inesperado. A equipe de desenvolvimento foi notificada.",
	}
	if err := ectx.JSON(http.StatusInternalServerError, response); err != nil {
		slog.Error("Failed to write internal server error response", "error", err)
	}
}

func formatValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Este campo é obrigatório."
	case "email":
		return "O valor fornecido não é um email válido."
	case "min":
		return fmt.Sprintf("Este campo deve ter no mínimo %s caracteres.", err.Param())
	case "max":
		return fmt.Sprintf("Este campo deve ter no máximo %s caracteres.", err.Param())
	case "url":
		return "O valor fornecido não é uma URL válida."
	default:
		return fmt.Sprintf("A validação para a regra '%s' falhou.", err.Tag())
	}
}
