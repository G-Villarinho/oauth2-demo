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

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) error
	GetByClientID(ctx context.Context, clientID string) (*entities.Client, error)
}

type clientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) ClientRepository {
	return &clientRepository{
		collection: db.Collection("clients"),
	}
}

func (r *clientRepository) Create(ctx context.Context, client *entities.Client) error {
	if client.ID.IsZero() {
		client.ID = primitive.NewObjectID()
	}

	if client.CreatedAt.IsZero() {
		client.CreatedAt = time.Now().UTC()
	}

	if _, err := r.collection.InsertOne(ctx, client); err != nil {
		return err
	}

	return nil
}

func (r *clientRepository) GetByClientID(ctx context.Context, clientID string) (*entities.Client, error) {
	var client entities.Client
	if err := r.collection.FindOne(ctx, bson.M{"client_id": clientID}).Decode(&client); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrClientNotFound
		}

		return nil, err
	}

	return &client, nil
}
