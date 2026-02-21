package httpapi

import "github.com/gin-demo/recipes-web/internal/module/recipe/domain"

func toDomainRequest(input createRecipeRequest) domain.CreateRecipeInput {
	return domain.CreateRecipeInput{
		Name:         input.Name,
		Tags:         input.Tags,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
	}
}

func toDomainUpdateRequest(id string, input updateRecipeRequest) domain.UpdateRecipeInput {
	return domain.UpdateRecipeInput{
		ID:           domain.RecipeID(id),
		Name:         input.Name,
		Tags:         input.Tags,
		Ingredients:  input.Ingredients,
		Instructions: input.Instructions,
	}
}
