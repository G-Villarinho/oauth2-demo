package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"github.com/aetheris-lab/aetheris-id/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateOTP(t *testing.T) {
	t.Run("should create OTP successfully when valid userID is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
		email := "test@example.com"
		userIDObj, _ := primitive.ObjectIDFromHex(userID)

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(otp *entities.OTP) bool {
				return otp.UserID == userIDObj
			})).
			Return(nil)

		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID, email)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userIDObj, result.UserID)
		assert.Equal(t, email, result.Email)
		assert.Len(t, result.Code, 6) // OTP deve ter 6 dígitos
		assert.True(t, result.ExpiresAt.After(time.Now().UTC()))
		assert.True(t, result.ExpiresAt.Before(time.Now().UTC().Add(11*time.Minute)))
	})

	t.Run("should return error when userID is invalid ObjectID format", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invalidUserID := "invalid-user-id"
		email := "test@example.com"

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, invalidUserID, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "convert userID to ObjectID")
	})

	t.Run("should return error when repository fails to create OTP", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
		email := "test@example.com"
		userIDObj, _ := primitive.ObjectIDFromHex(userID)

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(otp *entities.OTP) bool {
				return otp.UserID == userIDObj
			})).
			Return(errors.New("database connection failed"))

		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create otp")
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("should create OTP with correct expiration time based on config", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
		email := "test@example.com"
		userIDObj, _ := primitive.ObjectIDFromHex(userID)
		expirationMinutes := 5 * time.Minute

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(otp *entities.OTP) bool {
				return otp.UserID == userIDObj
			})).
			Return(nil)

		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: expirationMinutes,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID, email)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)

		expectedExpiration := time.Now().UTC().Add(expirationMinutes)
		// Permitir uma pequena diferença de tempo (1 segundo) devido ao tempo de execução
		assert.True(t, result.ExpiresAt.After(expectedExpiration.Add(-time.Second)))
		assert.True(t, result.ExpiresAt.Before(expectedExpiration.Add(time.Second)))
	})

	t.Run("should generate OTP code with correct length and format", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
		email := "test@example.com"
		userIDObj, _ := primitive.ObjectIDFromHex(userID)

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(otp *entities.OTP) bool {
				return otp.UserID == userIDObj
			})).
			Return(nil)

		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID, email)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Code, 6)

		// Verificar se o código contém apenas dígitos
		for _, char := range result.Code {
			assert.True(t, char >= '0' && char <= '9', "OTP code should contain only digits")
		}
	})

	t.Run("should return error when userID is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		emptyUserID := ""
		email := "test@example.com"

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, emptyUserID, email)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "convert userID to ObjectID")
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
		emptyEmail := ""
		userIDObj, _ := primitive.ObjectIDFromHex(userID)

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().
			Create(ctx, mock.MatchedBy(func(otp *entities.OTP) bool {
				return otp.UserID == userIDObj && otp.Email == emptyEmail
			})).
			Return(nil)

		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, &config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID, emptyEmail)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userIDObj, result.UserID)
		assert.Equal(t, emptyEmail, result.Email)
		assert.Len(t, result.Code, 6)
	})
}

func TestValidateCode(t *testing.T) {
	t.Run("should validate code and delete OTP successfully", func(t *testing.T) {
		ctx := context.Background()
		otpID := "otp-id"
		code := "123456"
		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().FindByID(ctx, otpID).Return(otp, nil)
		mockOTPRepo.EXPECT().Delete(ctx, otpID).Return(nil)

		config := configs.Environment{}
		otpService := NewOTPService(mockOTPRepo, &config)

		result, err := otpService.ValidateCode(ctx, code, otpID)

		require.NoError(t, err)
		assert.Equal(t, otp, result)
	})

	t.Run("should return error when OTP not found", func(t *testing.T) {
		ctx := context.Background()
		otpID := "otp-id"
		code := "123456"

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		errNotFound := errors.New("not found")
		mockOTPRepo.EXPECT().FindByID(ctx, otpID).Return(nil, errNotFound)

		config := configs.Environment{}
		otpService := NewOTPService(mockOTPRepo, &config)

		result, err := otpService.ValidateCode(ctx, code, otpID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "find otp by id")
	})

	t.Run("should return error when code is invalid", func(t *testing.T) {
		ctx := context.Background()
		otpID := "otp-id"
		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			Code:      "654321",
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().FindByID(ctx, otpID).Return(otp, nil)

		config := configs.Environment{}
		otpService := NewOTPService(mockOTPRepo, &config)

		result, err := otpService.ValidateCode(ctx, "123456", otpID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate code")
	})

	t.Run("should return error when code is expired", func(t *testing.T) {
		ctx := context.Background()
		otpID := "otp-id"
		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().FindByID(ctx, otpID).Return(otp, nil)

		config := configs.Environment{}
		otpService := NewOTPService(mockOTPRepo, &config)

		result, err := otpService.ValidateCode(ctx, "123456", otpID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validate code")
	})

	t.Run("should return error when delete fails after validation", func(t *testing.T) {
		ctx := context.Background()
		otpID := "otp-id"
		code := "123456"
		otp := &entities.OTP{
			ID:        primitive.NewObjectID(),
			Code:      code,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		mockOTPRepo.EXPECT().FindByID(ctx, otpID).Return(otp, nil)
		mockOTPRepo.EXPECT().Delete(ctx, otpID).Return(errors.New("delete error"))

		config := configs.Environment{}
		otpService := NewOTPService(mockOTPRepo, &config)

		result, err := otpService.ValidateCode(ctx, code, otpID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "delete otp")
	})
}
