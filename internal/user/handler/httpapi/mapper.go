package httpapi

import "github.com/gin-demo/recipes-web/internal/user/domain"

func toUpdateUserRequest(req updateUserRequest) domain.UpdateUserInput {
	return domain.UpdateUserInput{
		UserName: req.UserName,
		Password: req.Password,
	}
}

func toUserID(req reqUserID) domain.UserID {
	return domain.UserID(req.ID)
}
