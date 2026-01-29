package model

import "time"

// RecipeID represents a unique identifier for recipes.
type RecipeID string

// Recipe represents a cooking recipe with ingredients and instructions.
type Recipe struct {
	// ID is the unique identifier for the recipe
	ID RecipeID `json:"id" bson:"_id"`
	// Name is the title of the recipe
	Name string `json:"name" bson:"name"`
	// Tags is a list of tags associated with the recipe
	Tags []string `json:"tags" bson:"tags"`
	// Ingredients is a list of ingredients needed for the recipe
	Ingredients []string `json:"ingredients" bson:"ingredients"`
	// Instructions is a list of steps to prepare the recipe
	Instructions []string `json:"instructions" bson:"instructions"`
	// PublishedAt is the timestamp when the recipe was published
	PublishedAt time.Time `json:"publishedAt" bson:"publishedAt"`
}
