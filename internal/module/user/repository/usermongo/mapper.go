package usermongo

import (
	"github.com/gin-demo/recipes-web/internal/module/user/domain"
)

func ToDomainUser(doc userDocument) domain.User {
	return domain.User{
		ID:           domain.UserID(doc.ID),
		UserName:     doc.UserName,
		PasswordHash: doc.PasswordHash,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func FromDomainUser(user domain.User) (userDocument, error) {

	return userDocument{
		ID:           string(user.ID),
		UserName:     user.UserName,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
