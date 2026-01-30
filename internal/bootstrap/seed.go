package bootstrap

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/repository/mongorepo"
)

// SeedRecipe populates the repository with initial recipe data from the specified file.
func SeedRecipe(ctx context.Context, repo *mongorepo.Repository, seedPath string) error {
	return repo.SeedFromFile(ctx, seedPath)
}
