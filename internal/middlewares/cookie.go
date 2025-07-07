package middlewares

import (
	"net/http"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/labstack/echo/v4"
)

type CookieMiddleware interface {
	SetCookie(ectx echo.Context, value string, maxAge int)
	GetCookie(ectx echo.Context) (string, error)
	DeleteCookie(ectx echo.Context)
}

type cookieMiddleware struct {
	config *configs.Environment
}

func NewCookieMiddleware(config *configs.Environment) CookieMiddleware {
	return &cookieMiddleware{
		config: config,
	}
}

func (m *cookieMiddleware) SetCookie(ectx echo.Context, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     m.config.Security.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   m.config.Env == configs.Production,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	}

	ectx.SetCookie(cookie)
}

func (m *cookieMiddleware) GetCookie(ectx echo.Context) (string, error) {
	cookie, err := ectx.Cookie(m.config.Security.CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (m *cookieMiddleware) DeleteCookie(ectx echo.Context) {
	cookie := &http.Cookie{
		Name:     m.config.Security.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   m.config.Env == configs.Production,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   0,
	}

	ectx.SetCookie(cookie)
}
