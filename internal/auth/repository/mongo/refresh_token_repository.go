package mongo

import (
	"context"
	"errors"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	authColl *mongo.Collection
}

func New(coll *mongo.Collection) *Repository {
	return &Repository{authColl: coll}
}

func (repo *Repository) Save(ctx context.Context, token *domain.RefreshToken) error {
	_, err := repo.authColl.InsertOne(ctx, fromDomain(token))
	if err != nil {
		return domain.ErrPersistence
	}

	return nil
}

func (repo *Repository) FindByTokenHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var doc refreshTokenDoc
	if err := repo.authColl.FindOne(ctx, bson.M{"tokenHash": hash}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrPersistence
	}

	return toDomain(doc), nil
}

func (repo *Repository) DeleteByID(ctx context.Context, id domain.TokenID) error {
	update := bson.M{
		"$set": bson.M{
			"revoked": true,
		},
	}

	result, err := repo.authColl.UpdateOne(ctx, bson.M{"_id": string(id)}, update)
	if err != nil {
		return domain.ErrPersistence
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
