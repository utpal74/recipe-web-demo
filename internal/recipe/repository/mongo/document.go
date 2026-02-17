package mongo

import (
	"time"
)

// Recipe represents a cooking recipe with ingredients and instructions.
type recipeDocument struct {
	ID           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Tags         []string  `bson:"tags"`
	Ingredients  []string  `bson:"ingredients"`
	Instructions []string  `bson:"instructions"`
	PublishedAt  time.Time `bson:"publishedAt"`
}
