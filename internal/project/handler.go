package project

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *ProjectService
}

func NewProjectHandler(svc *ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: svc,
	}
}

func (ph *ProjectHandler) HandleGetAllProjects(c *gin.Context) {

	data, err := ph.service.GetAllProjects(c.Request.Context())
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting project data", data))
}

func (ph *ProjectHandler) PostCreateProject(c *gin.Context) {

	creatorId, _ := middleware.GetUserID(c)
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, insErr := ph.service.CreateProject(c.Request.Context(), req, creatorId)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successsfully creating new project", data))

}

func (ph *ProjectHandler) DeleteProjectHandle(c *gin.Context) {

	param := c.Param("id")
	projId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id on parameter! please provide valid uuid"))
		return
	}

	// handle delete chat room sama chat media nanti kalo udah dibikin

	delErr := ph.service.DeleteProjects(c.Request.Context(), projId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully remove the projects", nil))
}

func (ph *ProjectHandler) RegisterRoutes(r gin.IRouter) {

	projectApi := r.Group("/project")
	projectApi.Use(middleware.AuthMiddleware())

	projectAdmin := projectApi.Group("/")
	projectAdmin.Use(middleware.RoleMiddleware(middleware.RoleAdmin))
	projectAdmin.GET("/get-all", ph.HandleGetAllProjects)
	projectAdmin.POST("/add", ph.PostCreateProject)
	projectAdmin.DELETE("/delete/:id", ph.DeleteProjectHandle)
}
