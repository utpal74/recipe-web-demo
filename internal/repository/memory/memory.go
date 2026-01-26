package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/gin-demo/recipes-web/model"
	"github.com/rs/xid"
)

var (
	ErrNotFound      = errors.New("recipe not found")
	ErrPersistence   = errors.New("persistence failure")
	ErrIOFailure     = errors.New("IO failure")
	ErrSerialization = errors.New("serialization/deserialziation failure")
)

type Repository struct {
	mu       sync.RWMutex
	data     []model.Recipe
	dataPath string
}

func New(path string) (*Repository, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIOFailure, err)
	}
	defer file.Close()

	var recipes []model.Recipe
	if err := json.NewDecoder(file).Decode(&recipes); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSerialization, err)
	}

	return &Repository{data: recipes, dataPath: path}, nil
}

func (repo *Repository) Create(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	newRecipe := model.Recipe{
		ID:           model.RecipeID(xid.New().String()),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  time.Now(),
	}

	repo.data = append(repo.data, newRecipe)

	if err := saveAll(repo.dataPath, repo.data); err != nil {
		return model.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
	}

	return newRecipe, nil
}

func (repo *Repository) GetByID(ctx context.Context, id model.RecipeID) (model.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, recipe := range repo.data {
		if recipe.ID == id {
			return recipe, nil
		}
	}

	return model.Recipe{}, ErrNotFound
}

func (repo *Repository) GetAll(ctx context.Context) ([]model.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	out := make([]model.Recipe, len(repo.data))
	copy(out, repo.data)
	return out, nil
}

func (repo *Repository) Update(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i, r := range repo.data {
		select {
		case <-ctx.Done():
			return model.Recipe{}, ctx.Err()
		default:
		}

		if r.ID == recipe.ID {
			updated := append([]model.Recipe(nil), repo.data...)
			updated[i] = recipe

			if err := saveAll(repo.dataPath, updated); err != nil {
				return model.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
			}

			repo.data = updated
			return recipe, nil
		}
	}

	return model.Recipe{}, ErrNotFound
}

func (repo *Repository) Delete(ctx context.Context, id model.RecipeID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i, r := range repo.data {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if r.ID == id {
			updated := append(repo.data[0:i], repo.data[i+1:]...)

			if err := saveAll(repo.dataPath, updated); err != nil {
				return fmt.Errorf("%w: %v", ErrPersistence, err)
			}

			repo.data = updated
			return nil
		}
	}

	return ErrNotFound
}

func (repo *Repository) GetByTag(ctx context.Context, tag string) ([]model.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	recipes := []model.Recipe{}
	for _, r := range repo.data {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if slices.Contains(r.Tags, tag) {
			recipes = append(recipes, r)
		}
	}

	return recipes, nil
}

func saveAll(path string, recipes []model.Recipe) error {
	bytes, err := json.MarshalIndent(&recipes, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
