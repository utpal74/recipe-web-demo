package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRecipeJSONMarshalUnmarshal(t *testing.T) {
	now := time.Now()
	recipe := Recipe{
		ID:           "test-id",
		Name:         "Test Recipe",
		Tags:         []string{"tag1", "tag2"},
		Ingredients:  []string{"ing1", "ing2"},
		Instructions: []string{"step1", "step2"},
		PublishedAt:  now,
	}

	// Marshal
	data, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal recipe: %v", err)
	}

	// Unmarshal
	var unmarshaled Recipe
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal recipe: %v", err)
	}

	// Check fields
	if unmarshaled.ID != recipe.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, recipe.ID)
	}
	if unmarshaled.Name != recipe.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, recipe.Name)
	}
	if len(unmarshaled.Tags) != len(recipe.Tags) {
		t.Errorf("Tags length mismatch: got %v, want %v", len(unmarshaled.Tags), len(recipe.Tags))
	}
	for i, tag := range recipe.Tags {
		if unmarshaled.Tags[i] != tag {
			t.Errorf("Tag %d mismatch: got %v, want %v", i, unmarshaled.Tags[i], tag)
		}
	}
	// Similarly for ingredients and instructions
	if len(unmarshaled.Ingredients) != len(recipe.Ingredients) {
		t.Errorf("Ingredients length mismatch")
	}
	for i, ing := range recipe.Ingredients {
		if unmarshaled.Ingredients[i] != ing {
			t.Errorf("Ingredient %d mismatch", i)
		}
	}
	if len(unmarshaled.Instructions) != len(recipe.Instructions) {
		t.Errorf("Instructions length mismatch")
	}
	for i, ins := range recipe.Instructions {
		if unmarshaled.Instructions[i] != ins {
			t.Errorf("Instruction %d mismatch", i)
		}
	}
	// PublishedAt might have precision issues, check approximate
	if unmarshaled.PublishedAt.Sub(recipe.PublishedAt) > time.Second {
		t.Errorf("PublishedAt mismatch")
	}
}

func TestRecipeIDString(t *testing.T) {
	id := RecipeID("test")
	if string(id) != "test" {
		t.Errorf("RecipeID string conversion failed")
	}
}