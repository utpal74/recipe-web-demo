package mongorepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-demo/recipes-web/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (repo *Repository) SeedFromFile(ctx context.Context, filePath string) error {
	if repo.dbName == "" {
		return errors.New("MONGO_DATABASE is not set")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading recipe.json: %w", err)
	}

	var recipes = make([]model.Recipe, 0)
	if err := json.Unmarshal(data, &recipes); err != nil {
		return fmt.Errorf("error unmarshalling: %w", err)
	}

	collection := repo.mongoclient.Database(repo.dbName).Collection("recipes")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("recipes already exists, skip insert operation")
		return nil
	}

	docs := make([]interface{}, len(recipes))
	for i, r := range recipes {
		docs[i] = r
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("error inserting records in db: %w", err)
	}

	log.Printf("%d records inserted in DB", len(result.InsertedIDs))
	return nil
}
