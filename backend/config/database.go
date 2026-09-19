package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongo(cfg *Config) (*mongo.Client, error) {
	client, err := mongo.Connect(
		options.Client().ApplyURI(cfg.MongoURI),
	)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
