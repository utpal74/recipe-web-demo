package usermongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/user/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Repository provides MongoDB-backed user persistence.
type Repository struct {
	userColl *mongo.Collection
}

func New(collection *mongo.Collection) *Repository {
	// New creates a new Repository for user persistence using MongoDB.
	return &Repository{
		userColl: collection,
	}
}

// Create implements domain.Repository.
func (r *Repository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	newUserDoc, err := FromDomainUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %v", domain.ErrPersistence, err)
	}

	filteredUser := r.userColl.FindOne(ctx, bson.M{
		"userName": user.UserName,
	})

	if filteredUser.Acknowledged {
		return domain.User{}, domain.ErrUserAlreadyExists
	}

	result, err := r.userColl.InsertOne(ctx, newUserDoc)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %v", domain.ErrPersistence, err)
	}

	objID, ok := result.InsertedID.(string)
	if !ok {
		return domain.User{}, fmt.Errorf("%w: %v", domain.ErrPersistence, err)
	}

	newUserDoc.ID = objID

	return ToDomainUser(newUserDoc), nil
}

// Delete implements domain.Repository.
func (r *Repository) Delete(ctx context.Context, id domain.UserID) error {
	filter := bson.M{"id": id}

	_, err := r.userColl.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// FindByID implements domain.Repository.
func (r *Repository) FindByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	filter := bson.M{"id": id}

	userDoc := userDocument{}

	err := r.userColl.FindOne(ctx, filter).Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return domain.User{}, nil
}

// FindByUserName implements domain.Repository.
func (r *Repository) FindByUserName(ctx context.Context, userName string) (domain.User, error) {
	filter := bson.M{"userName": userName}

	var userDoc userDocument
	err := r.userColl.FindOne(ctx, filter).Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return ToDomainUser(userDoc), nil
}

// Update implements domain.Repository.
func (r *Repository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	filter := bson.M{"_id": user.ID}

	setFields := bson.M{}

	if user.UserName != "" {
		setFields["userName"] = user.UserName
	}

	if user.PasswordHash != "" {
		setFields["passwordHash"] = user.PasswordHash
	}

	setFields["updatedAt"] = time.Now()

	// Prevent empty update
	if len(setFields) == 0 {
		return domain.User{}, domain.ErrNothingToUpdate
	}

	update := bson.M{
		"$set": setFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated domain.User
	err := r.userColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		return domain.User{}, domain.ErrPersistence
	}

	return updated, nil
}
