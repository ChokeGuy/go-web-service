package database

import (
	"context"
	"fmt"
	"time"
	"web-service/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDBClient() (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s/?maxPoolSize=%s&w=majority", config.Env.DBHost, config.Env.DBPort, config.Env.DBPoolSize)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure client options
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	// Create new client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, client.Disconnect(ctx)
	}

	fmt.Printf("Connect to database %s successfully on PORT %s\n", config.Env.DBName, config.Env.DBPort)
	return client, nil
}

func DisconnectMongoDB(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
