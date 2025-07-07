package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// GetFullName retorna o nome completo do usuário
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidEmail verifica se o email é válido
func (u *User) IsValidEmail() bool {
	return u.Email != "" && len(u.Email) > 3 && len(u.Email) < 255
}
