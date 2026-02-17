package main

import (
	"context"
	"log"
	"os"
	"time"

	recipemongo "github.com/gin-demo/recipes-web/internal/recipe/repository/mongo"
	"github.com/gin-demo/recipes-web/internal/shared/bootstrap"
)

func main() {
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    mr, err := bootstrap.InitMongo(uri, "recipe-app")
    if err != nil {
        log.Fatalf("failed to init mongo: %v", err)
    }
    defer mr.Client.Disconnect(ctx)

    repo := recipemongo.New(mr.Database.Collection("recipes"))

    if err := bootstrap.SeedRecipe(ctx, repo, mr.Database.Collection("recipes"), "data/recipe.json"); err != nil {
        log.Fatalf("seeding failed: %v", err)
    }

    log.Println("seeding complete")
}
