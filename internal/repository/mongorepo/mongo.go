package mongorepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/model"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const RECIPE_COLLECTION = "recipes"

// Repository implements the recipe repository interface using MongoDB.
type Repository struct {
	mongoclient *mongo.Client
	dbName      string
}

// New creates a new Repository instance connected to the specified MongoDB URI and database.
func New(uri string, dbName string) (*Repository, error) {
	if uri == "" || dbName == "" {
		return nil, fmt.Errorf("%w: %v", domain.ErrPersistence, "MONGO_URI or MONGO_DATABASE can't be empty")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrPersistence, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(
		ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrPersistence, err)
	}

	log.Println("Connected to Mongo DB !!!")
	return &Repository{mongoclient: client, dbName: dbName}, nil
}

// Create adds a new recipe to the repository.
func (repo *Repository) Create(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	newRecipe := model.Recipe{
		ID:           model.RecipeID(xid.New().String()),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  time.Now(),
	}

	collection := repo.collection(RECIPE_COLLECTION)
	_, err := collection.InsertOne(ctx, newRecipe)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					return model.Recipe{}, fmt.Errorf("%w", domain.ErrConflict)
				}
			}
		}
		return model.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	return newRecipe, nil
}

// GetByID retrieves a recipe by its ID.
func (repo *Repository) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	collection := repo.collection(RECIPE_COLLECTION)
	var recipe model.Recipe

	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&recipe)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Recipe{}, domain.ErrNotFound
		}
		return model.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	return recipe, nil
}

// GetAll returns all recipes in the repository.
func (repo *Repository) GetAll(ctx context.Context) ([]model.Recipe, error) {
	collection := repo.collection(RECIPE_COLLECTION)

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}
	defer cur.Close(ctx)

	recipes := make([]model.Recipe, 0)
	for cur.Next(ctx) {
		var recipe model.Recipe
		if err := cur.Decode(&recipe); err != nil {
			return nil, domain.ErrPersistence
		}
		recipes = append(recipes, recipe)
	}

	if err := cur.Err(); err != nil {
		return nil, domain.ErrPersistence
	}

	return recipes, nil
}

// Update modifies an existing recipe in the repository.
func (repo *Repository) Update(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	collection := repo.collection(RECIPE_COLLECTION)

	filter := bson.M{"_id": recipe.ID}
	update := bson.M{
		"$set": bson.M{
			"name":         recipe.Name,
			"tags":         recipe.Tags,
			"ingredients":  recipe.Ingredients,
			"instructions": recipe.Instructions,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	var updated model.Recipe
	err = collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		return model.Recipe{}, domain.ErrPersistence
	}

	return updated, nil
}

// Delete removes a recipe from the repository by ID.
func (repo *Repository) Delete(ctx context.Context, id model.RecipeID) error {
	collection := repo.collection(RECIPE_COLLECTION)

	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%w", domain.ErrPersistence)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w", domain.ErrNotFound)
	}

	return nil
}

// GetByTag retrieves all recipes that contain the specified tag.
func (repo *Repository) GetByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	collection := repo.collection(RECIPE_COLLECTION)

	filter := bson.M{"tags": tag}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrPersistence
	}
	defer cur.Close(ctx)

	recipes := make([]model.Recipe, 0)
	for cur.Next(ctx) {
		var r model.Recipe
		if err := cur.Decode(&r); err != nil {
			return nil, domain.ErrPersistence
		}
		recipes = append(recipes, r)
	}

	if err := cur.Err(); err != nil {
		return nil, domain.ErrPersistence
	}

	if len(recipes) == 0 {
		return nil, domain.ErrNotFound
	}

	return recipes, nil
}

func (repo *Repository) collection(name string) *mongo.Collection {
	return repo.mongoclient.Database(repo.dbName).Collection(name)
}

// Close disconnects the MongoDB client.
func (repo *Repository) Close(ctx context.Context) error {
	return repo.mongoclient.Disconnect(ctx)
}
