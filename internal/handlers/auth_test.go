package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Run("should return success when valid email is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		otpToken := "test-otp-token"
		expiresAt := time.Now().Add(10 * time.Minute)

		expectedResponse := &models.SendVerificationCodeResponse{
			OTPToken:  otpToken,
			ExpiresAt: expiresAt,
		}

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			SendVerificationCode(ctx, email).
			Return(expectedResponse, nil)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)
		mockCookieMiddleware.EXPECT().
			SetCookie(mock.AnythingOfType("*echo.context"), otpToken, mock.AnythingOfType("int")).
			Return()

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: email}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.SendVerificationCodeResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, otpToken, response.OTPToken)
		assert.Equal(t, expiresAt.Unix(), response.ExpiresAt.Unix())
	})

	t.Run("should return bad request when payload binding fails", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context with invalid JSON
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrBadRequest, err)
	})

	t.Run("should return bad request when email validation fails", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: "invalid-email"}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.Error(t, err)
		assert.NotEqual(t, echo.ErrBadRequest, err) // Should be validation error, not bad request
	})

	t.Run("should return bad request when email is empty", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: ""}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.Error(t, err)
		assert.NotEqual(t, echo.ErrBadRequest, err) // Should be validation error, not bad request
	})

	t.Run("should return not found when user is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			SendVerificationCode(ctx, email).
			Return(nil, domain.ErrUserNotFound)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: email}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrNotFound, err)
	})

	t.Run("should return internal server error when auth service fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			SendVerificationCode(ctx, email).
			Return(nil, errors.New("database connection failed"))

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: email}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrInternalServerError, err)
	})

	t.Run("should set cookie with correct max age when successful", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		otpToken := "test-otp-token"
		expiresAt := time.Now().Add(10 * time.Minute)

		expectedResponse := &models.SendVerificationCodeResponse{
			OTPToken:  otpToken,
			ExpiresAt: expiresAt,
		}

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			SendVerificationCode(ctx, email).
			Return(expectedResponse, nil)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)
		mockCookieMiddleware.EXPECT().
			SetCookie(mock.AnythingOfType("*echo.context"), otpToken, mock.MatchedBy(func(maxAge int) bool {
				// Verifica se o maxAge está próximo do esperado (com tolerância de 1 segundo)
				expectedMaxAge := int(expiresAt.Sub(time.Now().UTC()).Seconds())
				return maxAge >= expectedMaxAge-1 && maxAge <= expectedMaxAge+1
			})).
			Return()

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.LoginPayload{Email: email}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Login(ectx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("should return success when valid code and OTP ID are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"
		accessToken := "test-access-token"
		expiresAt := time.Now().Add(10 * time.Minute)

		expectedResponse := &models.AuthenticateResponse{
			AccessToken: accessToken,
			ExpiresAt:   expiresAt,
		}

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(expectedResponse, nil)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)
		mockCookieMiddleware.EXPECT().
			SetCookie(mock.AnythingOfType("*echo.context"), accessToken, mock.AnythingOfType("int")).
			Return()

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.AuthenticateResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, accessToken, response.AccessToken)
		assert.Equal(t, expiresAt.Unix(), response.ExpiresAt.Unix())
	})

	t.Run("should return bad request when payload binding fails", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context with invalid JSON
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrBadRequest, err)
	})

	t.Run("should return bad request when code validation fails", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: ""} // Empty code
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.NotEqual(t, echo.ErrBadRequest, err) // Should be validation error, not bad request
	})

	t.Run("should return unauthorized when OTP ID is not found", func(t *testing.T) {
		// Arrange
		mockAuthService := mocks.NewAuthServiceMock(t)
		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: "123456"}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Don't set OTP claims in context

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})

	t.Run("should return unauthorized when OTP is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "nonexistent-otp-id"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(nil, domain.ErrOTPNotFound)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})

	t.Run("should return unauthorized when code is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "invalid-code"
		otpID := "test-otp-id"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(nil, domain.ErrInvalidCode)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})

	t.Run("should return unauthorized when OTP is expired", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "expired-otp-id"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(nil, domain.ErrOTPExpired)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrUnauthorized, err)
	})

	t.Run("should return internal server error when auth service fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(nil, errors.New("database connection failed"))

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, echo.ErrInternalServerError, err)
	})

	t.Run("should set cookie with correct max age when successful", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"
		accessToken := "test-access-token"
		expiresAt := time.Now().Add(10 * time.Minute)

		expectedResponse := &models.AuthenticateResponse{
			AccessToken: accessToken,
			ExpiresAt:   expiresAt,
		}

		mockAuthService := mocks.NewAuthServiceMock(t)
		mockAuthService.EXPECT().
			Authenticate(ctx, code, otpID).
			Return(expectedResponse, nil)

		mockCookieMiddleware := mocks.NewCookieMiddlewareMock(t)
		mockCookieMiddleware.EXPECT().
			SetCookie(mock.AnythingOfType("*echo.context"), accessToken, mock.MatchedBy(func(maxAge int) bool {
				// Verifica se o maxAge está próximo do esperado (com tolerância de 1 segundo)
				expectedMaxAge := int(expiresAt.Sub(time.Now().UTC()).Seconds())
				return maxAge >= expectedMaxAge-1 && maxAge <= expectedMaxAge+1
			})).
			Return()

		handler := NewAuthHandler(mockAuthService, mockCookieMiddleware)

		// Setup Echo context
		e := echo.New()
		e.Validator = &customValidator{validator: validator.New()}

		payload := models.AuthenticatePayload{Code: code}
		jsonPayload, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ectx := e.NewContext(req, rec)

		// Set OTP claims in context
		otpClaims := &models.OTPTokenClaims{}
		otpClaims.ID = otpID
		ectx.Set("otp", otpClaims)

		// Act
		err := handler.Authenticate(ectx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
