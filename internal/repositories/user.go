package repositories

import (
	"context"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (u *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	filter := bson.M{"email": email}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	filter := bson.M{"_id": id}
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
