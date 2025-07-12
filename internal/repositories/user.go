package repositories

import (
	"context"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain"
	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (u *userRepository) Create(ctx context.Context, user *entities.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	user.CreatedAt = time.Now()
	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	filter := bson.M{"email": email}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidObjectID
	}

	filter := bson.M{"_id": objectID}
	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
