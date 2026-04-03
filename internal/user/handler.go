package user

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(svc *UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}

func (uh *UserHandler) HandleGetAllUser(c *gin.Context) {

	data, err := uh.service.GetAllUser(c.Request.Context())
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully gathering user data", data))
}

func (uh *UserHandler) HandleSearchUser(c *gin.Context) {
	query := c.Query("q")
	data, err := uh.service.SearchUser(c.Request.Context(), query)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully gathering user data", data))
}

func (uh *UserHandler) HandleMyProfile(c *gin.Context) {
	userId, _ := middleware.GetUserID(c)

	userData, err := uh.service.GetUserById(c.Request.Context(), userId)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully gathering your profile data", userData))
}

func (uh *UserHandler) UpdateCurrentUser(c *gin.Context) {
	var updatedData UpdateUserDto
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid request body!"))
		return
	}
	myId, _ := middleware.GetUserID(c)
	data, updatedErr := uh.service.UpdateUserData(c.Request.Context(), updatedData, myId)
	if updatedErr != nil {
		c.JSON(updatedErr.Code, updatedErr)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "succesfully updated your data", gin.H{
		"updated_at": data,
	}))
}

func (uh *UserHandler) UpdateUserHandler(c *gin.Context) {
	var updatedData UpdateUserDto
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid request body!"))
		return
	}

	userIdStr := c.Param("id")

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid request parameter!"))
		return
	}
	data, updatedErr := uh.service.UpdateUserData(c.Request.Context(), updatedData, userId)
	if updatedErr != nil {
		c.JSON(updatedErr.Code, updatedErr)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "succesfully updated your data", gin.H{
		"updated_at": data,
	}))
}

func (uh *UserHandler) DeleteUserHandler(c *gin.Context) {

	param := c.Param("id")
	userId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, "invalid id parameter!")
		return
	}

	delErr := uh.service.DeleteUser(c.Request.Context(), userId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully deleted user!", gin.H{
		"deleted_at": time.Now(),
	}))
}

func (uh *UserHandler) HandleGetById(c *gin.Context) {

	param := c.Param("id")

	userId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id parameter!"))
		return
	}
	data, getErr := uh.service.GetUserById(c.Request.Context(), userId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting the user data", data))

}

func (uh *UserHandler) RegisterRoutes(r gin.IRouter) {

	userApi := r.Group("/user")
	userApi.Use(middleware.AuthMiddleware())
	userApi.GET("/my-profile", uh.HandleMyProfile)
	userApi.PATCH("/update-my-profile", uh.UpdateCurrentUser)

	// admin only
	userAdminOnly := userApi.Group("/")
	userAdminOnly.Use(middleware.RoleMiddleware(middleware.RoleAdmin))

	userAdminOnly.GET("/get-all", uh.HandleGetAllUser)
	userAdminOnly.PATCH("/update/:id", uh.UpdateUserHandler)
	userAdminOnly.GET("/search", uh.HandleSearchUser)
	userAdminOnly.DELETE("/delete/:id", uh.DeleteUserHandler)
	userAdminOnly.GET("/id/:id", uh.HandleGetById)
}
