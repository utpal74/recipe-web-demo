package httpapi

import (
	"github.com/gin-demo/recipes-web/internal/auth/usecase"
)

func RequestToUsecase(input signInUserRequest) usecase.SignInInput {
	return usecase.SignInInput{
		UserName: input.UserName,
		Password: input.Password,
	}
}

func UsecaseToResponse(input usecase.SignInOutput) signInUserResponse {
	return signInUserResponse{
		Token:   input.AccessToken,
		Expires: input.Expires,
	}
}

func toCreateUserRequest(req createUserRequest) usecase.SignUpInput {
	return usecase.SignUpInput{
		UserName: req.UserName,
		Password: req.Password,
	}
}

func toUserResponse(userName string) userResponse {
	return userResponse{
		UserName: userName,
	}
}
