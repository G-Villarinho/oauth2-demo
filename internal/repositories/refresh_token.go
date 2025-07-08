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

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error)
}

type refreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) RefreshTokenRepository {
	return &refreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	if refreshToken.ID.IsZero() {
		refreshToken.ID = primitive.NewObjectID()
	}

	refreshToken.CreatedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	filter := bson.M{"token_hash": tokenHash}

	var refreshToken entities.RefreshToken
	err := r.collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRefreshTokenNotFound
		}

		return nil, err
	}

	return &refreshToken, nil
}
