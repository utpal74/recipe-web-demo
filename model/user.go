package model

import "time"

// UserID presents unique ID for a user.
type UserID string

// User defines a user.
type User struct {
	ID        UserID    `json:"id" bson:"_id"`
	UserName  string    `json:"userName" bson:"userName"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

/*
TODO:
		*************************************
******* user belongs to Domain not model  *
********************************************
--------------------------------------------------
How to know if something belongs in domain (simple test)

Ask:

“If I removed HTTP, DB, cache — would this concept still exist?”

If yes → domain
If no → transport / infra / model / DTO

User passes this test easily.
-----------------------------------------

Later move this user.go to domain folder and remove json and bson tag because domain shouldn't know about json or monog specific thing. It should look like below at domain level.

type User struct {
	ID UserID
	UserName string
	CreatedAt time.Time
}


Create a mapper later at mongo layer such as

type UserDocument struct {
	ID string `bson:"_id"`
	UserName string `bson:"userName"`
	PasswordHash string `bson:"password"`
}
*/
