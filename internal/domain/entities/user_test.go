package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser_GetFullName(t *testing.T) {
	t.Run("should return full name when first and last name are provided", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		// Act
		fullName := user.GetFullName()

		// Assert
		assert.Equal(t, "João Silva", fullName)
	})

	t.Run("should return only first name when last name is empty", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		// Act
		fullName := user.GetFullName()

		// Assert
		assert.Equal(t, "João ", fullName)
	})

	t.Run("should return only last name when first name is empty", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "",
			LastName:  "Silva",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		// Act
		fullName := user.GetFullName()

		// Assert
		assert.Equal(t, " Silva", fullName)
	})
}

func TestUser_IsValidEmail(t *testing.T) {
	t.Run("should return true when email is valid", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@example.com",
			CreatedAt: time.Now(),
		}

		// Act
		isValid := user.IsValidEmail()

		// Assert
		assert.True(t, isValid)
	})

	t.Run("should return false when email is empty", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "Silva",
			Email:     "",
			CreatedAt: time.Now(),
		}

		// Act
		isValid := user.IsValidEmail()

		// Assert
		assert.False(t, isValid)
	})

	t.Run("should return false when email is too short", func(t *testing.T) {
		// Arrange
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "Silva",
			Email:     "a@b",
			CreatedAt: time.Now(),
		}

		// Act
		isValid := user.IsValidEmail()

		// Assert
		assert.False(t, isValid)
	})

	t.Run("should return false when email is too long", func(t *testing.T) {
		// Arrange
		longEmail := strings.Repeat("a", 256) + "@example.com"
		user := &User{
			ID:        primitive.NewObjectID(),
			FirstName: "João",
			LastName:  "Silva",
			Email:     longEmail,
			CreatedAt: time.Now(),
		}

		// Act
		isValid := user.IsValidEmail()

		// Assert
		assert.False(t, isValid)
	})
}
