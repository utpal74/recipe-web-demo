package domain

import (
	"time"
)

type RecipeID string

type Recipe struct {
	ID           RecipeID
	Name         string
	Tags         []string
	Ingredients  []string
	Instructions []string
	PublishedAt  time.Time
}

type CreateRecipeInput struct {
	Name         string
	Tags         []string
	Ingredients  []string
	Instructions []string
}

type UpdateRecipeInput struct {
	ID           RecipeID
	Name         *string
	Tags         []string
	Ingredients  []string
	Instructions []string
}
