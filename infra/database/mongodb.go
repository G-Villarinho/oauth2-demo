package database

import (
	"context"
	"fmt"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(env *configs.Environment) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), env.MongoDB.ConnectionTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(env.MongoDB.ConnectionURI).
		SetMaxPoolSize(env.MongoDB.MaxPoolSize).
		SetMaxConnIdleTime(env.MongoDB.MaxConnIdleTime)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo database: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo database: %w", err)
	}

	return client, nil
}

func NewMongoDatabase(client *mongo.Client, env *configs.Environment) *mongo.Database {
	return client.Database(env.MongoDB.DatabaseName)
}
