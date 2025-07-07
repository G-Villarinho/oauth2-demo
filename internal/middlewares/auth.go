package middlewares

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware interface {
	EnsureAuthenticated() echo.MiddlewareFunc
	EnsureOTPAuthenticated() echo.MiddlewareFunc
	AttachUserClaimsIfAuthenticated() echo.MiddlewareFunc
}

type authMiddleware struct {
	jwtService       services.JWTService
	cookieMiddleware CookieMiddleware
}

func NewAuthMiddleware(
	jwtService services.JWTService,
	cookieMiddleware CookieMiddleware,
) AuthMiddleware {
	return &authMiddleware{
		jwtService:       jwtService,
		cookieMiddleware: cookieMiddleware,
	}
}

func (m *authMiddleware) EnsureAuthenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			token, err := m.cookieMiddleware.GetCookie(ectx)
			if err != nil {
				return echo.ErrUnauthorized
			}

			claims, err := m.jwtService.ValidateAccessTokenJWT(ectx.Request().Context(), token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			SetUserClaims(ectx, &claims)

			return next(ectx)
		}
	}
}

func (m *authMiddleware) EnsureOTPAuthenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			token, err := m.cookieMiddleware.GetCookie(ectx)
			if err != nil {
				return echo.ErrUnauthorized
			}

			claims, err := m.jwtService.ValidateOTPTokenJWT(ectx.Request().Context(), token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			SetOTPClaims(ectx, &claims)

			return next(ectx)
		}
	}
}

func (m *authMiddleware) AttachUserClaimsIfAuthenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			token, err := m.cookieMiddleware.GetCookie(ectx)
			if err != nil {
				return next(ectx)
			}

			claims, err := m.jwtService.ValidateAccessTokenJWT(ectx.Request().Context(), token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			SetUserClaims(ectx, &claims)

			return next(ectx)
		}
	}
}
