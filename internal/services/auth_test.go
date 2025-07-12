package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendVerificationCode(t *testing.T) {
	t.Run("should return success when valid email is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		otpID := primitive.NewObjectID()
		expiresAt := time.Now().Add(10 * time.Minute)
		expectedToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."

		user := &entities.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     email,
			CreatedAt: time.Now(),
		}

		otp := &entities.OTP{
			ID:        otpID,
			UserID:    userID,
			Code:      "123456",
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			CreateOTP(ctx, userID.Hex(), email).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateOTPTokenJWT(ctx, otpID.Hex(), expiresAt).
			Return(expectedToken, nil)

		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedToken, result.OTPToken)
		assert.Equal(t, expiresAt, result.ExpiresAt)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, domain.ErrUserNotFound)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
		assert.Contains(t, err.Error(), domain.ErrUserNotFound.Error())
	})

	t.Run("should return error when user repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		expectedError := errors.New("database connection failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, expectedError)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when OTP creation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		expectedError := errors.New("OTP creation failed")

		user := &entities.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     email,
			CreatedAt: time.Now(),
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			CreateOTP(ctx, userID.Hex(), email).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when JWT generation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		otpID := primitive.NewObjectID()
		expiresAt := time.Now().Add(10 * time.Minute)
		expectedError := errors.New("JWT generation failed")

		user := &entities.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     email,
			CreatedAt: time.Now(),
		}

		otp := &entities.OTP{
			ID:        otpID,
			UserID:    userID,
			Code:      "123456",
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(user, nil)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			CreateOTP(ctx, userID.Hex(), email).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateOTPTokenJWT(ctx, otpID.Hex(), expiresAt).
			Return("", expectedError)

		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "generate otp token jwt")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := ""

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, domain.ErrUserNotFound)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
		assert.Contains(t, err.Error(), domain.ErrUserNotFound.Error())
	})

	t.Run("should return error when context is cancelled", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately
		email := "test@example.com"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, context.Canceled)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
		assert.Contains(t, err.Error(), context.Canceled.Error())
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("should return success when valid code and OTP ID are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"
		userID := primitive.NewObjectID()
		expectedToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."

		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateAccessTokenJWT(ctx, userID.Hex(), mock.AnythingOfType("time.Time")).
			Return(expectedToken, nil)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedToken, result.AccessToken)
		assert.True(t, result.ExpiresAt.After(time.Now()))
		assert.True(t, result.ExpiresAt.Before(time.Now().Add(11*time.Minute)))
	})

	t.Run("should return error when OTP validation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "invalid-code"
		otpID := "test-otp-id"
		expectedError := domain.ErrInvalidCode

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when OTP is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "nonexistent-otp-id"
		expectedError := domain.ErrOTPNotFound

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when OTP is expired", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "expired-otp-id"
		expectedError := domain.ErrOTPExpired

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when JWT generation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"
		userID := primitive.NewObjectID()
		expectedError := errors.New("JWT generation failed")

		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateAccessTokenJWT(ctx, userID.Hex(), mock.AnythingOfType("time.Time")).
			Return("", expectedError)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "generate access token jwt")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when code is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := ""
		otpID := "test-otp-id"
		expectedError := domain.ErrInvalidCode

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when OTP ID is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := ""
		expectedError := domain.ErrOTPNotFound

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when context is cancelled", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately
		code := "123456"
		otpID := "test-otp-id"
		expectedError := context.Canceled

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should use correct expiration time from config", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "test-otp-id"
		userID := primitive.NewObjectID()
		expectedToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9..."
		configExpiration := 15 * time.Minute

		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: configExpiration,
			},
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateAccessTokenJWT(ctx, userID.Hex(), mock.AnythingOfType("time.Time")).
			Return(expectedToken, nil)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, &config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedToken, result.AccessToken)
		assert.True(t, result.ExpiresAt.After(time.Now()))
		assert.True(t, result.ExpiresAt.Before(time.Now().Add(configExpiration+time.Minute)))
	})
}
