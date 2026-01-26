package main

import (
	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/internal/handler/httpapi"
	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-gonic/gin"
)

const DATA_PATH = "data/recipe.json"

func main() {
	/*
		GET /recipes - Return list of recipes
		GET /recipes/{id} - Get recipe by ID
		POST /recipes - Create new recipe
		PUT /recipes/{id} - Updates an existing recipes
		DELETE /recipes/{id} - Deletes an existing recipes
		GET /recipes/search?tag=X = Search recipe by tag
	*/

	router := gin.Default()

	repo, err := memory.New(DATA_PATH)
	if err != nil {
		panic(err)
	}
	ctrl := recipe.New(repo)
	handler := httpapi.New(ctrl)

	router.GET("/recipes", handler.ListRecipeHandler)
	router.GET("/recipes/search", handler.ListRecipesByTagHandler)
	router.GET("/recipes/:id", handler.GetRecipeByIDHandler)
	router.POST("/recipes", handler.CreateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)

	router.Run()
}
