package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestAuthorize(t *testing.T) {
	const userID = "test-user-id"
	const userIDKey = "user_id"

	t.Run("should return success and redirect when valid payload is provided", func(t *testing.T) {
		// Arrange
		payload := models.AuthorizePayload{
			ClientID:            "test-client-id",
			RedirectURI:         "http://localhost/callback",
			ResponseType:        "code",
			Scope:               "openid",
			State:               "xyz",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
		}
		jsonPayload, _ := json.Marshal(payload)

		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(userIDKey, userID)

		mockOAuthService := mocks.NewOAuthServiceMock(t)
		handler := NewOAuthHandler(mockOAuthService)

		expectedResponse := &models.AuthorizeResponse{
			RedirectURL: "http://localhost/callback?code=123456&state=xyz",
		}

		mockOAuthService.EXPECT().
			Authorize(mock.Anything, mock.AnythingOfType("models.AuthorizeInput")).
			Return(expectedResponse, nil).Once()

		// Act
		err := handler.Authorize(c)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, expectedResponse.RedirectURL, rec.Header().Get(echo.HeaderLocation))
	})

	t.Run("should return bad request when binding fails", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuthService := mocks.NewOAuthServiceMock(t)
		handler := NewOAuthHandler(mockOAuthService)

		// Act
		err := handler.Authorize(c)

		// Assert
		require.Error(t, err)
		httpErr, ok := err.(*echo.HTTPError)
		require.True(t, ok)
		assert.Equal(t, echo.ErrBadRequest.Code, httpErr.Code)
	})

	t.Run("should return bad request when validation fails", func(t *testing.T) {
		// Arrange
		payload := models.AuthorizePayload{} // Invalid payload
		jsonPayload, _ := json.Marshal(payload)

		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockOAuthService := mocks.NewOAuthServiceMock(t)
		handler := NewOAuthHandler(mockOAuthService)

		// Act
		err := handler.Authorize(c)

		// Assert
		require.Error(t, err)
	})

	testCases := []struct {
		name          string
		serviceError  error
		expectedError *echo.HTTPError
	}{
		{"should return bad request when client is not found", domain.ErrClientNotFound, echo.ErrBadRequest},
		{"should return bad request when redirect uri is invalid", domain.ErrInvalidRedirectURI, echo.ErrBadRequest},
		{"should return bad request when grant type is invalid", domain.ErrInvalidGrantType, echo.ErrBadRequest},
		{"should return bad request when response type is invalid", domain.ErrInvalidResponseType, echo.ErrBadRequest},
		{"should return bad request when scope is invalid", domain.ErrInvalidScope, echo.ErrBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			payload := models.AuthorizePayload{
				ClientID:            "test-client-id",
				RedirectURI:         "http://localhost/callback",
				ResponseType:        "code",
				Scope:               "openid",
				State:               "xyz",
				CodeChallenge:       "test-challenge",
				CodeChallengeMethod: "S256",
			}
			jsonPayload, _ := json.Marshal(payload)

			e := echo.New()
			e.Validator = &customValidator{validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set(userIDKey, userID)

			mockOAuthService := mocks.NewOAuthServiceMock(t)
			handler := NewOAuthHandler(mockOAuthService)

			mockOAuthService.EXPECT().
				Authorize(mock.Anything, mock.AnythingOfType("models.AuthorizeInput")).
				Return(nil, tc.serviceError).Once()

			// Act
			err := handler.Authorize(c)

			// Assert
			require.Error(t, err)
			httpErr, ok := err.(*echo.HTTPError)
			require.True(t, ok)
			assert.Equal(t, tc.expectedError.Code, httpErr.Code)
		})
	}

	t.Run("should return internal server error for other service errors", func(t *testing.T) {
		// Arrange
		payload := models.AuthorizePayload{
			ClientID:            "test-client-id",
			RedirectURI:         "http://localhost/callback",
			ResponseType:        "code",
			Scope:               "openid",
			State:               "xyz",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: "S256",
		}
		jsonPayload, _ := json.Marshal(payload)

		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(userIDKey, userID)

		mockOAuthService := mocks.NewOAuthServiceMock(t)
		handler := NewOAuthHandler(mockOAuthService)

		expectedErr := errors.New("some unexpected error")
		mockOAuthService.EXPECT().
			Authorize(mock.Anything, mock.AnythingOfType("models.AuthorizeInput")).
			Return(nil, expectedErr).Once()

		// Act
		err := handler.Authorize(c)

		// Assert
		require.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}
