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

	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/rs/xid"
)

var (
	ErrNotFound      = errors.New("recipe not found")
	ErrPersistence   = errors.New("persistence failure")
	ErrIOFailure     = errors.New("IO failure")
	ErrSerialization = errors.New("serialization/deserialziation failure")
)

// Repository implements the recipe repository interface using in-memory storage with file persistence.
// Repository provides in-memory recipe persistence.
type Repository struct {
	mu       sync.RWMutex
	data     []domain.Recipe
	dataPath string
}

// New creates a new Repository instance with data loaded from the specified file path.
func New(path string) (*Repository, error) {
	// New creates a new Repository for recipe persistence using a file path.
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIOFailure, err)
	}
	defer file.Close()

	var recipes []domain.Recipe
	if err := json.NewDecoder(file).Decode(&recipes); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSerialization, err)
	}

	return &Repository{data: recipes, dataPath: path}, nil
}

// Create adds a new recipe to the repository.
func (repo *Repository) Create(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	newRecipe := domain.Recipe{
		ID:           domain.RecipeID(xid.New().String()),
		Name:         recipe.Name,
		Tags:         recipe.Tags,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		PublishedAt:  time.Now(),
	}

	repo.data = append(repo.data, newRecipe)

	if err := saveAll(repo.dataPath, repo.data); err != nil {
		return domain.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
	}

	return newRecipe, nil
}

func (repo *Repository) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	return nil
}

// GetByID retrieves a recipe by its ID.
func (repo *Repository) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, recipe := range repo.data {
		if recipe.ID == id {
			return recipe, nil
		}
	}

	return domain.Recipe{}, ErrNotFound
}

// GetAll returns all recipes in the repository.
func (repo *Repository) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	out := make([]domain.Recipe, len(repo.data))
	copy(out, repo.data)
	return out, nil
}

// Update modifies an existing recipe in the repository.
func (repo *Repository) Update(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i, r := range repo.data {
		select {
		case <-ctx.Done():
			return domain.Recipe{}, ctx.Err()
		default:
		}

		if r.ID == recipe.ID {
			updated := append([]domain.Recipe(nil), repo.data...)
			updated[i] = recipe

			if err := saveAll(repo.dataPath, updated); err != nil {
				return domain.Recipe{}, fmt.Errorf("%w: %v", ErrPersistence, err)
			}

			repo.data = updated
			return recipe, nil
		}
	}

	return domain.Recipe{}, ErrNotFound
}

// Delete removes a recipe from the repository by ID.
func (repo *Repository) Delete(ctx context.Context, id domain.RecipeID) error {
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

// GetByTag retrieves all recipes that contain the specified tag.
func (repo *Repository) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	recipes := []domain.Recipe{}
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

func saveAll(path string, recipes []domain.Recipe) error {
	bytes, err := json.MarshalIndent(&recipes, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}
