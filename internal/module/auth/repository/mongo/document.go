package mongo

import "time"

type refreshTokenDoc struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"userID"`
	TokenHash string    `bson:"tokenHash"`
	Revoked   bool      `bson:"revoked"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresAt time.Time `bson:"expiresAt"`
}
