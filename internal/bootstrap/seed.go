package bootstrap

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/repository/mongorepo"
)

func SeedRecipe(ctx context.Context, repo *mongorepo.Repository, seedPath string) error {
	return repo.SeedFromFile(ctx, seedPath)
}
