package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-demo/recipes-web/internal/bootstrap"
	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/internal/handler/httpapi"
	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-demo/recipes-web/internal/repository/mongorepo"
	"github.com/gin-gonic/gin"
)

var (
	repo      recipe.RecipeRepository
	mongoRepo *mongorepo.Repository
	err       error
)

type Config struct {
	RepoType string
	DataPath string
	MongoURI string
	HttpAddr string
	SeedData bool
}

func main() {
	/*
		GET /recipes - Return list of recipes
		GET /recipes/{id} - Get recipe by ID
		POST /recipes - Create new recipe
		PUT /recipes/{id} - Updates an existing recipes
		DELETE /recipes/{id} - Deletes an existing recipes
		GET /recipes/search?tag=X = Search recipe by tag
	*/

	cfg := loadConfig()

	switch cfg.RepoType {
	case "memory":
		repo, err = memory.New(cfg.DataPath)
	case "mongo":
		mongoRepo, err := mongorepo.New(cfg.MongoURI, "recipes")
		if err != nil {
			log.Fatalf("failed to initialize mongo repository: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if cfg.SeedData {
			if err := bootstrap.SeedRecipe(ctx, mongoRepo, cfg.DataPath); err != nil {
				log.Fatal(err)
			}
		}

		repo = mongoRepo

	default:
		log.Fatalf("unknown REPO_TYPE: %s", cfg.RepoType)
	}

	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	router := gin.Default()

	ctrl := recipe.New(repo)
	handler := httpapi.New(ctrl)

	router.GET("/recipes", handler.ListRecipeHandler)
	router.GET("/recipes/search", handler.ListRecipesByTagHandler)
	router.GET("/recipes/:id", handler.GetRecipeByIDHandler)
	router.POST("/recipes", handler.CreateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)

	srv := &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: router,
	}

	go func() {
		log.Printf("HTTP server listening on %s", cfg.HttpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if mongoRepo != nil {
		log.Println("Closing MongoDB connection...")
		if err := mongoRepo.Close(ctx); err != nil {
			log.Printf("Mongo close error: %v", err)
		}
	}

	log.Println("Server exiting")
}

func loadConfig() Config {
	// default configuration
	cfg := Config{
		RepoType: "memory",
		DataPath: "data/recipe.json",
		MongoURI: "mongodb://admin:password@localhost:27017/test?authSource=admin",
		HttpAddr: ":8080",
	}

	if v := os.Getenv("REPO_TYPE"); v != "" {
		cfg.RepoType = v
	}
	if v := os.Getenv("DATA_PATH"); v != "" {
		cfg.DataPath = v
	}
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.MongoURI = v
	}
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		cfg.HttpAddr = v
	}
	if v := os.Getenv("SEED_DATA"); v != "" {
		value, err := strconv.ParseBool(v)
		if err != nil {
			fmt.Printf("error parsing SEED_DATA env variable: %v\n", err)
		}
		cfg.SeedData = value
	}

	return cfg
}
