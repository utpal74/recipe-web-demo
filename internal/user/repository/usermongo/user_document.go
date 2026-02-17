package usermongo

import (
	"time"
)

type userDocument struct {
	ID           string    `bson:"_id"`
	UserName     string    `bson:"userName"`
	PasswordHash string    `bson:"passwordHash"`
	Role         string    `bson:"role"`
	CreatedAt    time.Time `bson:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt"`
}
