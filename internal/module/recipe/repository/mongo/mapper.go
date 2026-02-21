package mongo

import (
	"github.com/gin-demo/recipes-web/internal/module/recipe/domain"
)

func toDomainRecipe(doc recipeDocument) domain.Recipe {
	return domain.Recipe{
		ID:           domain.RecipeID(doc.ID),
		Name:         doc.Name,
		Tags:         doc.Tags,
		Ingredients:  doc.Ingredients,
		Instructions: doc.Instructions,
		PublishedAt:  doc.PublishedAt,
	}
}

func fromDomainRecipe(recipe domain.Recipe) recipeDocument {
	return recipeDocument{
		ID:           string(recipe.ID),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  recipe.PublishedAt,
	}
}
