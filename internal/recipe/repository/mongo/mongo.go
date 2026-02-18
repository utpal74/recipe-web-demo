package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Repository implements the recipe repository interface using MongoDB.
// Repository provides MongoDB-backed recipe persistence.
type Repository struct {
	recipeColl *mongo.Collection
}

// New creates a new Repository instance connected to the specified MongoDB URI and database.
func New(collection *mongo.Collection) *Repository {
	// New creates a new Repository for recipe persistence using MongoDB.
	return &Repository{recipeColl: collection}
}

// Create adds a new recipe to the repository.
func (repo *Repository) Create(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	_, err := repo.recipeColl.InsertOne(ctx, fromDomainRecipe(recipe))
	if err != nil {
		var writeErr mongo.WriteException

		if errors.As(err, &writeErr) {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					return domain.Recipe{}, fmt.Errorf("%w", domain.ErrConflict)
				}
			}
		}

		return domain.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	return recipe, nil
}

func (repo *Repository) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	docs := make([]recipeDocument, len(recipes))

	for i, r := range recipes {
		docs[i] = fromDomainRecipe(r)
	}

	_, err := repo.recipeColl.InsertMany(ctx, docs)
	var bulkErr mongo.BulkWriteException
	if err != nil {
		if errors.As(err, &bulkErr) {
			for _, e := range bulkErr.WriteErrors {
				if e.Code == 11000 {
					return domain.ErrConflict
				}
			}
		}

		return domain.ErrPersistence
	}

	return nil
}

// GetByID retrieves a recipe by its ID.
func (repo *Repository) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	var doc recipeDocument
	filter := bson.M{"_id": string(id)}

	err := repo.recipeColl.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Recipe{}, domain.ErrNotFound
		}
		return domain.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	return toDomainRecipe(doc), nil
}

// GetAll returns all recipes in the repository.
func (repo *Repository) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	cur, err := repo.recipeColl.Find(ctx, bson.M{})
	if err != nil {
		return []domain.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	defer cur.Close(ctx)

	recipes := make([]domain.Recipe, 0)
	for cur.Next(ctx) {
		var doc recipeDocument
		if err := cur.Decode(&doc); err != nil {
			return nil, domain.ErrPersistence
		}
		recipes = append(recipes, toDomainRecipe(doc))
	}

	if err := cur.Err(); err != nil {
		return nil, domain.ErrPersistence
	}

	return recipes, nil
}

// Update modifies an existing recipe in the repository.
func (repo *Repository) Update(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	filter := bson.M{"_id": string(recipe.ID)}
	update := bson.M{
		"$set": bson.M{
			"name":         recipe.Name,
			"tags":         recipe.Tags,
			"ingredients":  recipe.Ingredients,
			"instructions": recipe.Instructions,
		},
	}

	result, err := repo.recipeColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Recipe{}, fmt.Errorf("%w", domain.ErrPersistence)
	}

	if result.MatchedCount == 0 {
		return domain.Recipe{}, domain.ErrNotFound
	}

	var updated recipeDocument
	err = repo.recipeColl.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		return domain.Recipe{}, domain.ErrPersistence
	}

	return toDomainRecipe(updated), nil
}

// Delete removes a recipe from the repository by ID.
func (repo *Repository) Delete(ctx context.Context, id domain.RecipeID) error {
	filter := bson.M{"_id": string(id)}

	result, err := repo.recipeColl.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%w", domain.ErrPersistence)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%w", domain.ErrNotFound)
	}

	return nil
}

// GetByTag retrieves all recipes that contain the specified tag.
func (repo *Repository) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	filter := bson.M{"tags": tag}

	cur, err := repo.recipeColl.Find(ctx, filter)
	if err != nil {
		return nil, domain.ErrPersistence
	}
	defer cur.Close(ctx)

	recipes := make([]domain.Recipe, 0)
	for cur.Next(ctx) {
		var doc recipeDocument
		if err := cur.Decode(&doc); err != nil {
			return nil, domain.ErrPersistence
		}
		recipes = append(recipes, toDomainRecipe(doc))
	}

	if err := cur.Err(); err != nil {
		return nil, domain.ErrPersistence
	}

	return recipes, nil
}
