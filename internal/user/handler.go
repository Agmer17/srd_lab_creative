package user

import "github.com/gin-gonic/gin"

type UserHandler struct {
	service *UserService
}

func NewUserHandler(svc *UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}

func (uh *UserHandler) RegisterRoutes(r gin.IRouter) {

	// do something
}
