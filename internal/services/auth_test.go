package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendVerificationCode(t *testing.T) {
	t.Run("should return success when valid email is provided and all services work correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		otpID := primitive.NewObjectID()
		expiresAt := time.Now().Add(10 * time.Minute)
		expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

		user := &models.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     email,
			CreatedAt: time.Now(),
		}

		otp := &models.OTP{
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
			CreateOTP(ctx, userID.Hex()).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateOTPTokenJWT(ctx, otpID.Hex(), expiresAt).
			Return(expectedToken, nil)

		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedToken, result.OTPToken)
	})

	t.Run("should return error when user repository fails to find user by email", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"
		expectedError := errors.New("user not found")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, expectedError)

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("should return error when OTP service fails to create OTP", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		expectedError := errors.New("failed to create OTP")

		user := &models.User{
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
			CreateOTP(ctx, userID.Hex()).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create otp")
		assert.Contains(t, err.Error(), "failed to create OTP")
	})

	t.Run("should return error when JWT service fails to generate OTP token", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "test@example.com"
		userID := primitive.NewObjectID()
		otpID := primitive.NewObjectID()
		expiresAt := time.Now().Add(10 * time.Minute)
		expectedError := errors.New("failed to generate JWT")

		user := &models.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     email,
			CreatedAt: time.Now(),
		}

		otp := &models.OTP{
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
			CreateOTP(ctx, userID.Hex()).
			Return(otp, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateOTPTokenJWT(ctx, otpID.Hex(), expiresAt).
			Return("", expectedError)

		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "generate otp token jwt")
		assert.Contains(t, err.Error(), "failed to generate JWT")
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := ""

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByEmail(ctx, email).
			Return(nil, errors.New("invalid email"))

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)
		config := configs.Environment{}

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
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

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.SendVerificationCode(ctx, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by email")
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("should return success when valid code and otpID are provided and all services work correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "507f1f77bcf86cd799439011"
		userID := primitive.NewObjectID()
		expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

		otpIDObj, _ := primitive.ObjectIDFromHex(otpID)
		otp := &models.OTP{
			ID:        otpIDObj,
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		user := &models.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByID(ctx, userID.Hex()).
			Return(user, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateAccessTokenJWT(ctx, user.FirstName, user.LastName, user.Email, user.ID.Hex(), mock.AnythingOfType("time.Time")).
			Return(expectedToken, nil)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedToken, result.AccessToken)
	})

	t.Run("should return error when OTP service fails to validate code", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "invalid"
		otpID := "507f1f77bcf86cd799439011"
		expectedError := models.ErrInvalidCode

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), "invalid code")
	})

	t.Run("should return error when OTP code is expired", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "507f1f77bcf86cd799439011"
		expectedError := models.ErrOTPExpired

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
		assert.Contains(t, err.Error(), "otp expired")
	})

	t.Run("should return error when user repository fails to find user by ID", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "507f1f77bcf86cd799439011"
		userID := primitive.NewObjectID()
		expectedError := models.ErrUserNotFound

		otpIDObj, _ := primitive.ObjectIDFromHex(otpID)
		otp := &models.OTP{
			ID:        otpIDObj,
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

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByID(ctx, userID.Hex()).
			Return(nil, expectedError)

		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find user by id")
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("should return error when JWT service fails to generate access token", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := "507f1f77bcf86cd799439011"
		userID := primitive.NewObjectID()
		expectedError := errors.New("failed to generate JWT")

		otpIDObj, _ := primitive.ObjectIDFromHex(otpID)
		otp := &models.OTP{
			ID:        otpIDObj,
			UserID:    userID,
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		user := &models.User{
			ID:        userID,
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(otp, nil)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockUserRepo.EXPECT().
			FindByID(ctx, userID.Hex()).
			Return(user, nil)

		mockJWTService := mocks.NewJWTServiceMock(t)
		mockJWTService.EXPECT().
			GenerateAccessTokenJWT(ctx, user.FirstName, user.LastName, user.Email, user.ID.Hex(), mock.AnythingOfType("time.Time")).
			Return("", expectedError)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "generate access token jwt")
		assert.Contains(t, err.Error(), "failed to generate JWT")
	})

	t.Run("should return error when code is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := ""
		otpID := "507f1f77bcf86cd799439011"
		expectedError := models.ErrInvalidCode

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
	})

	t.Run("should return error when otpID is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		code := "123456"
		otpID := ""
		expectedError := models.ErrOTPNotFound

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, expectedError)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
	})

	t.Run("should return error when context is cancelled", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel the context immediately
		code := "123456"
		otpID := "507f1f77bcf86cd799439011"

		config := configs.Environment{
			OTP: configs.OTP{
				JWTExpirationMinutes: 10 * time.Minute,
			},
		}

		mockOTPService := mocks.NewOTPServiceMock(t)
		mockOTPService.EXPECT().
			ValidateCode(ctx, code, otpID).
			Return(nil, context.Canceled)

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockJWTService := mocks.NewJWTServiceMock(t)

		authService := NewAuthService(mockUserRepo, mockOTPService, mockJWTService, config)

		// Act
		result, err := authService.Authenticate(ctx, code, otpID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate otp")
	})
}
