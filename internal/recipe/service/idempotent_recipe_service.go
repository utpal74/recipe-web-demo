package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	authdomain "github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/gin-demo/recipes-web/internal/recipe/domain"
	"github.com/gin-demo/recipes-web/internal/shared/id"
	"github.com/gin-demo/recipes-web/internal/shared/idempotency"
)

type idempotentRecipeService struct {
	inner       domain.RecipeService
	store       idempotency.IdempotencyStore
	idGenerator id.ULIDGenerator
}

func NewIdempotentRecipeService(
	recipeService domain.RecipeService,
	store idempotency.IdempotencyStore,
	idGenerator id.ULIDGenerator,
) domain.RecipeService {
	return &idempotentRecipeService{
		inner:       recipeService,
		store:       store,
		idGenerator: idGenerator,
	}
}

// Create implements [domain.RecipeService].
// TODO:
// 1) Add idempotent header
// 2) remove hardcoded header key from middleware
func (i *idempotentRecipeService) Create(ctx context.Context, input domain.CreateRecipeInput) (domain.Recipe, error) {
	fmt.Printf("Locked ***** %v\n", "h1")
	identity, ok := ctx.Value("authIdentity").(authdomain.Identity)
	if !ok {
		return i.inner.Create(ctx, input)
	}

	idemKey, ok := ctx.Value("idempotencyKey").(string)
	if !ok || idemKey == "" {
		return i.inner.Create(ctx, input)
	}

	record := &idempotency.Record{
		ID:        i.idGenerator.New(),
		UserID:    identity.UserID,
		Key:       idemKey,
		Status:    idempotency.StatusProcessing,
		CreatedAt: time.Now(),
	}

	locked, err := i.store.TryLock(ctx, record)
	if err != nil {
		return domain.Recipe{}, err
	}

	if !locked {
		existing, err := i.store.Get(ctx, identity.UserID, idemKey)
		if err != nil {
			return domain.Recipe{}, err
		}

		if existing.Response != nil {
			var recipe domain.Recipe
			if err := json.Unmarshal(existing.Response, &recipe); err == nil {
				return recipe, nil
			}
		}
	}

	result, err := i.inner.Create(ctx, input)
	if err != nil {
		return result, err
	}

	bytes, err := json.Marshal(result)
	if err == nil {
		record.Response = bytes
		record.Status = idempotency.StatusCompleted
		record.StatusCode = 201

		_ = i.store.Set(ctx, record)
	}

	return result, nil
}

// CreateMany implements [domain.RecipeService].
func (i *idempotentRecipeService) CreateMany(ctx context.Context, recipes []domain.Recipe) error {
	return i.inner.CreateMany(ctx, recipes)
}

// Delete implements [domain.RecipeService].
func (i *idempotentRecipeService) Delete(ctx context.Context, id domain.RecipeID) error {
	return i.inner.Delete(ctx, id)
}

// GetAll implements [domain.RecipeService].
func (i *idempotentRecipeService) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	return i.inner.GetAll(ctx)
}

// GetByID implements [domain.RecipeService].
func (i *idempotentRecipeService) GetByID(ctx context.Context, id domain.RecipeID) (domain.Recipe, error) {
	return i.inner.GetByID(ctx, id)
}

// GetByTag implements [domain.RecipeService].
func (i *idempotentRecipeService) GetByTag(ctx context.Context, tag string) ([]domain.Recipe, error) {
	return i.inner.GetByTag(ctx, tag)
}

// Update implements [domain.RecipeService].
func (i *idempotentRecipeService) Update(ctx context.Context, input domain.UpdateRecipeInput) (domain.Recipe, error) {
	return i.inner.Update(ctx, input)
}
