package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "tokenHash", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.authColl.Indexes().CreateOne(ctx, index)
	return err
}
