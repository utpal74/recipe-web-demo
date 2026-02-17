package httpapi

type updateUserRequest struct {
	UserName *string `json:"userName,omitempty"`
	Password *string `json:"password,omitempty"`
}

type reqUserID struct {
	ID string `uri:"id" binding:"required"`
}

type reqUserName struct {
	Name string `uri:"name" binding:"required"`
}
