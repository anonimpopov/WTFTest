package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetMongoClient(connectURL string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(connectURL))
	err := client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
