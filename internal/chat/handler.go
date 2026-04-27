package chat

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	service *MessagingService
}

func NewChatHandler(svc *MessagingService) *ChatHandler {
	return &ChatHandler{
		service: svc,
	}
}

func (chh *ChatHandler) GetLatestChat(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)

	data, err := chh.service.GetLatestChat(c.Request.Context(), userId)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting the latest chat", data))
}

func (chh *ChatHandler) PostSendChat(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)
	param := c.Param("projectId")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	var req createChatDto
	if err := c.ShouldBind(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, insErr := chh.service.CreateProjectMessage(c.Request.Context(), userId, projectId, req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully creating a message", data))
}

func (chh *ChatHandler) RegisterRoutes(r gin.IRouter) {

	chatApi := r.Group("/chat")
	chatApi.Use(middleware.AuthMiddleware())

	chatApi.GET("/latest", chh.GetLatestChat)
	chatApi.POST("/group/:projectId/send", chh.PostSendChat)
}
