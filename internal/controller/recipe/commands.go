package recipe

// UpdateRecipeCommand contains the fields that can be updated for a recipe.
type UpdateRecipeCommand struct {
	// Name is the optional new name for the recipe
	Name *string
	// Tags is the optional new list of tags for the recipe
	Tags []string
	// Ingredients is the optional new list of ingredients for the recipe
	Ingredients []string
}
