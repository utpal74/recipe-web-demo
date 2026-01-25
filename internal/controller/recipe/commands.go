package recipe

type UpdateRecipeCommand struct {
	Name        *string
	Tags        []string
	Ingredients []string
}
