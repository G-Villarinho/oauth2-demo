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

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userIDObj, result.UserID)
		assert.Len(t, result.Code, 6) // OTP deve ter 6 dígitos
		assert.True(t, result.ExpiresAt.After(time.Now().UTC()))
		assert.True(t, result.ExpiresAt.Before(time.Now().UTC().Add(11*time.Minute)))
	})

	t.Run("should return error when userID is invalid ObjectID format", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		invalidUserID := "invalid-user-id"

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, invalidUserID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "convert userID to ObjectID")
	})

	t.Run("should return error when repository fails to create OTP", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "507f1f77bcf86cd799439011"
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

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID)

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

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID)

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

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, userID)

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

		mockOTPRepo := mocks.NewOTPRepositoryMock(t)
		config := configs.Environment{
			OTP: configs.OTP{
				ExpirationMinutes: 10 * time.Minute,
			},
		}

		otpService := NewOTPService(mockOTPRepo, config)

		// Act
		result, err := otpService.CreateOTP(ctx, emptyUserID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "convert userID to ObjectID")
	})
}
