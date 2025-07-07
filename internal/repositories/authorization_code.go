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

type AuthorizationCodeRepository interface {
	Create(ctx context.Context, authorizationCode *entities.AuthorizationCode) error
	FindByCode(ctx context.Context, code string) (*entities.AuthorizationCode, error)
	Delete(ctx context.Context, id string) error
}

type authorizationCodeRepository struct {
	collection *mongo.Collection
}

func NewAuthorizationCodeRepository(db *mongo.Database) AuthorizationCodeRepository {
	return &authorizationCodeRepository{
		collection: db.Collection("authorization_codes"),
	}
}

func (r *authorizationCodeRepository) Create(ctx context.Context, authorizationCode *entities.AuthorizationCode) error {
	if authorizationCode.ID.IsZero() {
		authorizationCode.ID = primitive.NewObjectID()
	}

	authorizationCode.CreatedAt = time.Now().UTC()

	if _, err := r.collection.InsertOne(ctx, authorizationCode); err != nil {
		return err
	}

	return nil
}

func (r *authorizationCodeRepository) FindByCode(ctx context.Context, code string) (*entities.AuthorizationCode, error) {
	filter := bson.M{"code": code}

	var authorizationCode entities.AuthorizationCode
	if err := r.collection.FindOne(ctx, filter).Decode(&authorizationCode); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrAuthorizationCodeNotFound
		}

		return nil, err
	}

	return &authorizationCode, nil
}

func (r *authorizationCodeRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if _, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID}); err != nil {
		return err
	}

	return nil
}
