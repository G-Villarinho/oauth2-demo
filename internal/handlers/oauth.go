package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/labstack/echo/v4"
)

type OAuthHandler interface {
	Authorize(ectx echo.Context) error
}

type oauthHandler struct {
	oauthService services.OAuthService
}

func NewOAuthHandler(oauthService services.OAuthService) OAuthHandler {
	return &oauthHandler{
		oauthService: oauthService,
	}
}

func (h *oauthHandler) Authorize(ectx echo.Context) error {
	logger := slog.With(
		slog.String("handler", "oauth"),
		slog.String("method", ectx.Request().Method),
		slog.String("path", ectx.Request().URL.Path),
	)

	var payload models.AuthorizePayload
	if err := ectx.Bind(&payload); err != nil {
		logger.Error("failed to bind input", "error", err)
		return echo.ErrBadRequest
	}

	if err := ectx.Validate(payload); err != nil {
		logger.Error("failed to validate payload", "error", err)
		return err
	}

	input := models.NewAuthorizeInput(payload, middlewares.GetUserID(ectx))

	response, err := h.oauthService.Authorize(ectx.Request().Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			logger.Warn(err.Error())
			return echo.ErrBadRequest
		}

		if errors.Is(err, domain.ErrInvalidRedirectURI) {
			logger.Warn(err.Error())
			return echo.ErrBadRequest
		}

		if errors.Is(err, domain.ErrInvalidGrantType) || errors.Is(err, domain.ErrInvalidResponseType) {
			logger.Warn(err.Error())
			return echo.ErrBadRequest
		}

		if errors.Is(err, domain.ErrInvalidScope) {
			logger.Warn(err.Error())
			return echo.ErrBadRequest
		}

		logger.Error("authorize", "error", err)
		return err
	}

	return ectx.Redirect(http.StatusFound, response.RedirectURL)
}
