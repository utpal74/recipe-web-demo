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

	authdomain "github.com/gin-demo/recipes-web/internal/module/auth/domain"
	authhandler "github.com/gin-demo/recipes-web/internal/module/auth/handler/httpapi"
	authmongo "github.com/gin-demo/recipes-web/internal/module/auth/repository/mongo"
	"github.com/gin-demo/recipes-web/internal/module/auth/service"
	recipedomain "github.com/gin-demo/recipes-web/internal/module/recipe/domain"
	"github.com/gin-demo/recipes-web/internal/module/recipe/handler/httpapi"
	"github.com/gin-demo/recipes-web/internal/module/recipe/repository"
	recipemongo "github.com/gin-demo/recipes-web/internal/module/recipe/repository/mongo"
	recipeservice "github.com/gin-demo/recipes-web/internal/module/recipe/service"
	userdomain "github.com/gin-demo/recipes-web/internal/module/user/domain"
	userhandler "github.com/gin-demo/recipes-web/internal/module/user/handler/httpapi"
	"github.com/gin-demo/recipes-web/internal/module/user/repository/usermongo"
	userservice "github.com/gin-demo/recipes-web/internal/module/user/service"
	"github.com/gin-demo/recipes-web/internal/platform/bootstrap"
	"github.com/gin-demo/recipes-web/internal/platform/cache/redisrecipe"
	"github.com/gin-demo/recipes-web/internal/platform/id"
	"github.com/gin-demo/recipes-web/internal/platform/idempotency"
	"github.com/gin-demo/recipes-web/internal/platform/middleware"
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
		recipeRepo       recipedomain.RecipeRepository
		userRepo         userdomain.UserRepository
		refreshTokenRepo authdomain.RefreshTokenRepository
		err              error
	)

	idGenerator := &id.ULIDGenerator{}
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found here, using system environment variable")
	}

	cfg := loadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoResource, err := bootstrap.InitMongo(cfg.MongoURI, "recipe-app")
	if err != nil {
		log.Fatalf("failed to initialize mongo repository: %v", err)
	}
	defer mongoResource.Client.Disconnect(ctx)

	recipeMongoRepo := recipemongo.New(mongoResource.Database.Collection("recipes"))

	userMongoRepo := usermongo.New(mongoResource.Database.Collection("users"))
	if err := userMongoRepo.EnsureIndexes(ctx); err != nil {
		log.Fatalf("error applying index: %v", err)
	}

	authMongoRepo := authmongo.New(mongoResource.Database.Collection("auth"))
	if err := authMongoRepo.EnsureIndexes(ctx); err != nil {
		log.Fatalf("error applying index: %v", err)
	}

	recipeRepo = recipeMongoRepo
	userRepo = userMongoRepo
	refreshTokenRepo = authMongoRepo

	if cfg.SeedData {
		if err := bootstrap.SeedRecipe(ctx, recipeRepo, mongoResource.Database.Collection("recipes"), cfg.DataPath); err != nil {
			log.Fatal(err)
		}
	}

	redisClient, err := bootstrap.NewRedis(os.Getenv("REDIS_ADDR"), "", 0)
	if err != nil {
		log.Printf("redis client init error : %v\n", err)
	}
	if redisClient != nil {
		cache := redisrecipe.NewCache(redisClient, 30*time.Minute)
		cachedRepo := repository.NewCachedRepository(recipeRepo, cache)
		recipeRepo = cachedRepo
	}

	router := gin.Default()

	store := idempotency.NewRedisStore(redisClient, 24*time.Hour)

	recipeService := recipeservice.NewRecipeService(recipeRepo, idGenerator)
	idemptentRecipeService := recipeservice.NewIdempotentRecipeService(
		recipeService,
		store,
		*idGenerator,
	)

	handler := httpapi.New(idemptentRecipeService)

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	tokenService := service.NewJwtTokenService(service.Config{
		Secret: os.Getenv("JWT_SECRET"),
		Issuer: "recipe-app",
	})

	jwtMiddleware := middleware.NewAuthMiddleWare(tokenService)

	pwdHasher := service.New()

	userService := userservice.NewUserService(userRepo, pwdHasher, idGenerator)

	config := service.AuthConfig{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}

	authService := service.NewAuthService(
		userService,
		tokenService,
		pwdHasher,
		refreshTokenRepo,
		idGenerator,
		config,
	)

	authHandler := authhandler.New(authService)

	userHandler := userhandler.New(userService)

	api := router.Group("/")

	// ----- Public Auth Endpoint -------
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authHandler.CreateUserHandler)
		auth.POST("/signin", authHandler.SignInHandler)
		auth.POST("/refresh", authHandler.RefreshHandler)
	}

	// ---- Protected Auth Endpoint --------
	authSecured := api.Group("/auth")
	authSecured.Use(jwtMiddleware.Handle())
	{
		authSecured.POST("/signout", authHandler.SignOutHandler)
	}

	// ----- Public Recipe Endpoint -------
	api.GET("/recipes", handler.ListRecipeHandler)
	api.GET("/recipes/search", handler.ListRecipesByTagHandler)

	// ---- Protected Endpoint --------

	secured := api.Group("/")
	secured.Use(jwtMiddleware.Handle())

	users := secured.Group("/users")
	{
		users.GET("/name/:name", userHandler.FindUserByNameHandler)
		users.GET("/:id", userHandler.FindUserByIDHandler)
		users.DELETE("/:id", userHandler.DeleteUserHandler)
		users.PUT("/:id", userHandler.UpdateUserHandler)
	}

	recipes := secured.Group("/recipes")
	{
		recipes.GET("/:id", handler.GetRecipeByIDHandler)
		recipes.POST("/", handler.CreateRecipeHandler)
		recipes.DELETE("/:id", handler.DeleteRecipeHandler)
		recipes.PUT("/:id", handler.UpdateRecipeHandler)
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

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
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
		RepoType: "mongo",
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
