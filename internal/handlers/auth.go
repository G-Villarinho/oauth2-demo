package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Login(ectx echo.Context) error
}

type authHandler struct {
	authService      services.AuthService
	cookieMiddleware middlewares.CookieMiddleware
}

func NewAuthHandler(
	authService services.AuthService,
	cookieMiddleware middlewares.CookieMiddleware,
) AuthHandler {
	return &authHandler{
		authService:      authService,
		cookieMiddleware: cookieMiddleware,
	}
}

func (h *authHandler) Login(ectx echo.Context) error {
	logger := slog.With(
		slog.String("handler", "auth"),
		slog.String("method", "login"),
	)

	var payload models.LoginPayload
	if err := ectx.Bind(&payload); err != nil {
		logger.Error("bind payload", "error", err)
		return echo.ErrBadRequest
	}

	if err := ectx.Validate(payload); err != nil {
		logger.Error("validate payload", "error", err)
		return err
	}

	response, err := h.authService.SendVerificationCode(ectx.Request().Context(), payload.Email)
	if err != nil {
		if errors.Is(err, entities.ErrUserNotFound) {
			logger.Error(err.Error())
			return echo.ErrNotFound
		}

		logger.Error("send verification code", "error", err)
		return echo.ErrInternalServerError
	}

	maxAge := int(response.ExpiresAt.Sub(time.Now().UTC()).Seconds())
	h.cookieMiddleware.SetCookie(ectx, response.OTPToken, maxAge)

	return ectx.JSON(http.StatusOK, response)
}
