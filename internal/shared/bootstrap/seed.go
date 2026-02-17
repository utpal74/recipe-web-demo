package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// recipeDTO is a local DTO used only for seeding purposes
type recipeDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	PublishedAt  string   `json:"publishedAt"`
}

// SeedRecipe populates the repository with initial recipe data from the specified file.
// If collection already contains data, seeding is skipped (idempotent behavior).
func SeedRecipe(ctx context.Context, recipeRepo domain.RecipeRepository, recipeColl *mongo.Collection, seedPath string) error {
	// Check if collection already has data
	count, err := recipeColl.EstimatedDocumentCount(ctx)
	if err != nil {
		return fmt.Errorf("error checking collection: %w", err)
	}

	if count > 0 {
		log.Printf("Collection already contains %d recipes, skipping seed", count)
		return nil
	}

	data, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("error reading recipe.json: %w", err)
	}

	var dtos []recipeDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return fmt.Errorf("error unmarshalling: %w", err)
	}

	// Convert DTOs to domain models
	recipes := make([]domain.Recipe, len(dtos))
	for i, dto := range dtos {
		publishedAt, err := time.Parse(time.RFC3339, dto.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}

		recipes[i] = domain.Recipe{
			ID:           domain.RecipeID(dto.ID),
			Name:         dto.Name,
			Tags:         dto.Tags,
			Ingredients:  dto.Ingredients,
			Instructions: dto.Instructions,
			PublishedAt:  publishedAt,
		}
	}

	log.Printf("Seeding %d recipes from %s", len(recipes), seedPath)
	return recipeRepo.CreateMany(ctx, recipes)
}
