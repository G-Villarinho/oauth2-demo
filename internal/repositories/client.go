package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByClientID(ctx context.Context, clientID string) (*models.Client, error)
}

type clientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) ClientRepository {
	return &clientRepository{
		collection: db.Collection("clients"),
	}
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	if client.ID.IsZero() {
		client.ID = primitive.NewObjectID()
	}

	if client.CreatedAt.IsZero() {
		client.CreatedAt = time.Now().UTC()
	}

	if _, err := r.collection.InsertOne(ctx, client); err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

func (r *clientRepository) GetByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	var client models.Client
	if err := r.collection.FindOne(ctx, bson.M{"client_id": clientID}).Decode(&client); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientNotFound
		}

		return nil, err
	}

	return &client, nil
}
