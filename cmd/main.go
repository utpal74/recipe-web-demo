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
	"github.com/gin-demo/recipes-web/internal/cache/redisrecipe"
	"github.com/gin-demo/recipes-web/internal/controller/recipe"
	"github.com/gin-demo/recipes-web/internal/domain"
	"github.com/gin-demo/recipes-web/internal/handler/httpapi"
	"github.com/gin-demo/recipes-web/internal/handler/httpapi/auth"
	"github.com/gin-demo/recipes-web/internal/handler/httpapi/middleware"
	"github.com/gin-demo/recipes-web/internal/repository"
	"github.com/gin-demo/recipes-web/internal/repository/memory"
	"github.com/gin-demo/recipes-web/internal/repository/mongorepo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config holds the application configuration from environment variables.
type Config struct {
	RepoType string
	DataPath string
	MongoURI string
	HttpAddr string
	SeedData bool
}

// main initializes and runs the recipe application server.
func main() {
	/*
		GET /recipes - Return list of recipes
		GET /recipes/{id} - Get recipe by ID
		POST /recipes - Create new recipe
		PUT /recipes/{id} - Updates an existing recipes
		DELETE /recipes/{id} - Deletes an existing recipes
		GET /recipes/search?tag=X = Search recipe by tag
	*/

	var (
		repo      domain.RecipeRepository
		mongoRepo *mongorepo.Repository
		err       error
	)

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found here, using system environment variable")
	}

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

	redisClient, err := bootstrap.NewRedis(os.Getenv("REDIS_ADDR"), "", 0)
	if err != nil {
		log.Printf("redis client init error : %v\n", err)
	}

	if redisClient != nil {
		cache := redisrecipe.NewCache(redisClient, 30*time.Minute)
		cachedRepo := repository.NewCachedRepository(repo, cache)
		repo = cachedRepo
	}

	router := gin.Default()

	ctrl := recipe.New(repo)
	handler := httpapi.New(ctrl)

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	authHandler := auth.New(auth.Config{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: "recipe-app",
	})

	router.POST("/signin", authHandler.SignInHandler)

	router.GET("/recipes", handler.ListRecipeHandler)
	router.GET("/recipes/search", handler.ListRecipesByTagHandler)

	authorized := router.Group("/recipes")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/:id", handler.GetRecipeByIDHandler)
		authorized.POST("/", handler.CreateRecipeHandler)
		authorized.DELETE("/:id", handler.DeleteRecipeHandler)
		authorized.PUT("/:id", handler.UpdateRecipeHandler)
	}

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

// loadConfig reads configuration from environment variables with defaults.
func loadConfig() Config {
	// default configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := Config{
		RepoType: "memory",
		DataPath: "data/recipe.json",
		MongoURI: os.Getenv("MONGO_URI"),
		HttpAddr: ":" + port,
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
