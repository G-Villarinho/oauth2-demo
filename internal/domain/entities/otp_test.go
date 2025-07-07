package entities

import (
	"testing"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOTP_IsExpired(t *testing.T) {
	t.Run("should return true when OTP is expired", func(t *testing.T) {
		// Arrange
		otp := &OTP{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			CreatedAt: time.Now(),
		}

		// Act
		isExpired := otp.IsExpired()

		// Assert
		assert.True(t, isExpired)
	})

	t.Run("should return false when OTP is not expired", func(t *testing.T) {
		// Arrange
		otp := &OTP{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(1 * time.Hour), // Expires in 1 hour
			CreatedAt: time.Now(),
		}

		// Act
		isExpired := otp.IsExpired()

		// Assert
		assert.False(t, isExpired)
	})
}

func TestOTP_ValidateCode(t *testing.T) {
	t.Run("should return nil when code is valid and not expired", func(t *testing.T) {
		// Arrange
		otp := &OTP{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(1 * time.Hour), // Expires in 1 hour
			CreatedAt: time.Now(),
		}

		// Act
		err := otp.ValidateCode("123456")

		// Assert
		require.NoError(t, err)
	})

	t.Run("should return error when code is invalid", func(t *testing.T) {
		// Arrange
		otp := &OTP{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(1 * time.Hour), // Expires in 1 hour
			CreatedAt: time.Now(),
		}

		// Act
		err := otp.ValidateCode("654321")

		// Assert
		require.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCode, err)
	})

	t.Run("should return error when OTP is expired", func(t *testing.T) {
		// Arrange
		otp := &OTP{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			Code:      "123456",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			CreatedAt: time.Now(),
		}

		// Act
		err := otp.ValidateCode("123456")

		// Assert
		require.Error(t, err)
		assert.Equal(t, domain.ErrOTPExpired, err)
	})
}
