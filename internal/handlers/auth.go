package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Login(ectx echo.Context) error
	Authenticate(ectx echo.Context) error
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
		if errors.Is(err, domain.ErrUserNotFound) {
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

func (h *authHandler) Authenticate(ectx echo.Context) error {
	logger := slog.With(
		slog.String("handler", "auth"),
		slog.String("method", "authenticate"),
	)

	var payload models.AuthenticatePayload
	if err := ectx.Bind(&payload); err != nil {
		logger.Error("bind payload", "error", err)
		return echo.ErrBadRequest
	}

	if err := ectx.Validate(payload); err != nil {
		logger.Error("validate payload", "error", err)
		return err
	}

	otpID := middlewares.GetOTPJTI(ectx)
	if otpID == "" {
		logger.Error("otp id not found")
		return echo.ErrUnauthorized
	}

	response, err := h.authService.Authenticate(ectx.Request().Context(), payload.Code, otpID)
	if err != nil {
		if errors.Is(err, domain.ErrOTPNotFound) {
			logger.Error(err.Error())
			return echo.ErrUnauthorized
		}

		if errors.Is(err, domain.ErrInvalidCode) {
			logger.Error(err.Error())
			return echo.ErrUnauthorized
		}

		if errors.Is(err, domain.ErrOTPExpired) {
			logger.Error(err.Error())
			return echo.ErrUnauthorized
		}

		logger.Error("authenticate", "error", err)
		return echo.ErrInternalServerError
	}

	maxAge := int(response.ExpiresAt.Sub(time.Now().UTC()).Seconds())
	h.cookieMiddleware.SetCookie(ectx, response.AccessToken, maxAge)

	return ectx.JSON(http.StatusOK, response)
}
