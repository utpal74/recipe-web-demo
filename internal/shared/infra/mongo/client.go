package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewClient(uri string) (*mongo.Client, error) {
	if uri == "" {
		return nil, errors.New("MONGO_URI can't be empty")
	}

	clientOpts, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := clientOpts.Ping(ctx,
		readpref.Primary()); err != nil {
		return nil, err
	}

	return clientOpts, nil
}
