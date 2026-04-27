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

func (chh *ChatHandler) GetChatDataFromProject(c *gin.Context) {
	userId, _ := middleware.GetUserID(c)

	param := c.Param("projectId")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	paramRoom := c.Param("room")
	roomId, err := uuid.Parse(paramRoom)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid room id"))
		return
	}

	data, getErr := chh.service.GetAllMessageFromProject(c.Request.Context(), userId, projectId, roomId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting chat data", data))
}

func (chh *ChatHandler) GetMediaAttachment(c *gin.Context) {
	param := c.Param("token")
	userId, _ := middleware.GetUserID(c)
	filename, err := chh.service.GetMediaAccessFromToken(c.Request.Context(), param, userId)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Header("X-Robots-Tag", "noindex, nofollow, noimageindex")
	c.Header("Cache-Control", "private, no-store")
	c.File(filename)
}

func (chh *ChatHandler) PostPersonalChat(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)

	param := c.Param("target")
	targetUuid, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	var req createPersonalChatDto
	if err := c.ShouldBind(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, insErr := chh.service.CreatePersonalChat(c.Request.Context(), targetUuid, userId, req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully sending the chat!", data))
}

func (chh *ChatHandler) DeleteChat(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)

	param := c.Param("id")
	targetUuid, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid chat id"))
		return
	}

	delErr := chh.service.DeleteChat(c.Request.Context(), targetUuid, userId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully delete chat data", nil))
}

func (chh *ChatHandler) GetPersonalChatData(c *gin.Context) {
	userId, _ := middleware.GetUserID(c)

	param := c.Param("roomId")
	roomId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid chat id"))
		return
	}

	data, getErr := chh.service.GetPersonalChatData(c.Request.Context(), roomId, userId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting chat data", data))
}

func (chh *ChatHandler) RegisterRoutes(r gin.IRouter) {

	chatApi := r.Group("/chat")
	chatApi.Use(middleware.AuthMiddleware())

	chatApi.GET("/latest", chh.GetLatestChat)
	chatApi.POST("/group/:projectId/send", chh.PostSendChat)
	chatApi.GET("/group/:projectId/:room", chh.GetChatDataFromProject)
	chatApi.POST("/personal/:target/send", chh.PostPersonalChat)
	chatApi.GET("/private-media/:token", chh.GetMediaAttachment)
	chatApi.GET("/personal/:roomId", chh.GetPersonalChatData)
	chatApi.DELETE("/delete/:id", chh.DeleteChat)
}
