package mongo

import (
	"github.com/gin-demo/recipes-web/internal/module/recipe/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toDomainRecipe(doc recipeDocument) domain.Recipe {
	var id string

	switch v := doc.ID.(type) {
	case primitive.ObjectID:
		id = v.Hex()
	case string:
		id = v
	default:
		id = ""
	}

	return domain.Recipe{
		ID:           domain.RecipeID(id),
		Name:         doc.Name,
		Tags:         doc.Tags,
		Ingredients:  doc.Ingredients,
		Instructions: doc.Instructions,
		PublishedAt:  doc.PublishedAt,
	}
}

func fromDomainRecipe(recipe domain.Recipe) recipeDocument {
	var id interface{}

	// try to convert to ObjectID
	if objID, err := primitive.ObjectIDFromHex(string(recipe.ID)); err == nil {
		id = objID // existing ObjectID
	} else {
		id = string(recipe.ID) // custom ID
	}

	return recipeDocument{
		ID:           id,
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  recipe.PublishedAt,
	}
}
