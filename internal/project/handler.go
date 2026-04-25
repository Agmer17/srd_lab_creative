package project

import (
	"fmt"

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

func (ph *ProjectHandler) PatchProjectData(c *gin.Context) {

	param := c.Param("id")
	projId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid uuid for project id! please provide a valid uuid"))
		return
	}

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, uptErr := ph.service.UpdateProjectData(c.Request.Context(), projId, req)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updating the project data", data))
}

func (ph *ProjectHandler) HandleGetDetail(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)
	param := c.Param("id")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	data, getErr := ph.service.GetDetailById(c.Request.Context(), projectId, userId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting the project data", data))
}

func (ph *ProjectHandler) HandleGetAllMember(c *gin.Context) {

	userId, _ := middleware.GetUserID(c)
	paramProj := c.Param("projectId")

	projectId, err := uuid.Parse(paramProj)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	data, getErr := ph.service.GetMemberFromProject(c.Request.Context(), projectId, userId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting the member data", data))
}

func (ph *ProjectHandler) PostNewMember(c *gin.Context) {
	creatorId, _ := middleware.GetUserID(c)

	projectParam := c.Param("projectId")

	projectId, err := uuid.Parse(projectParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	var req addNewMemberDto
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, addErr := ph.service.AddNewMember(c.Request.Context(), projectId, creatorId, req)
	if addErr != nil {
		c.JSON(addErr.Code, addErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully addning new member to project", data))
}

func (ph *ProjectHandler) PatchMemberData(c *gin.Context) {

	creatorId, _ := middleware.GetUserID(c)

	projectParam := c.Param("projectId")
	projectId, err := uuid.Parse(projectParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	var req updateMemberDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	memberId, err := uuid.Parse(req.MemberId)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid member id"))
		return
	}

	roleId, err := uuid.Parse(req.NewRole)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid role id"))
		return
	}

	data, uptErr := ph.service.UpdateProjectMemberRole(c.Request.Context(), creatorId, memberId, projectId, roleId, req.IsOwner)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updating the member data", data))

}

func (ph *ProjectHandler) DeleteProjectMember(c *gin.Context) {
	creatorId, _ := middleware.GetUserID(c)

	projectParam := c.Param("projectId")
	projectId, err := uuid.Parse(projectParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	delParam := c.Param("memberId")
	deletedId, err := uuid.Parse(delParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid member id"))
		return
	}

	delErr := ph.service.RemoveUserFromProject(c.Request.Context(), creatorId, deletedId, projectId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully removing user from this project", nil))

}

func (ph *ProjectHandler) GetProgressFromProject(c *gin.Context) {

	param := c.Param("projectId")

	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	userId, _ := middleware.GetUserID(c)

	data, getErr := ph.service.GetProgressFromProject(c.Request.Context(), projectId, userId)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting progress data", data))
}

func (ph *ProjectHandler) PostCreateProgress(c *gin.Context) {
	param := c.Param("projectId")

	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	userId, _ := middleware.GetUserID(c)
	var req createProgressRequests
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	newData, insErr := ph.service.CreateProjectProgress(c.Request.Context(), userId, projectId, req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return

	}

	c.JSON(200, shared.NewSuccessResponse(200, "success adding new progress to prject", newData))
}

func (ph *ProjectHandler) DeleteProjectProgress(c *gin.Context) {
	progParam := c.Param("id")
	progressId, err := uuid.Parse(progParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	delErr := ph.service.RemoveProjectProgress(c.Request.Context(), progressId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully deleting progress", nil))
}

func (ph *ProjectHandler) PatchUpdateProgress(c *gin.Context) {
	progParam := c.Param("id")
	progressId, err := uuid.Parse(progParam)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	var req updateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, uptErr := ph.service.UpdateProjectProgress(c.Request.Context(), progressId, req)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updating the project data", data))
}

func (ph *ProjectHandler) GetRevisionFromProject(c *gin.Context) {
	param := c.Param("projectId")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	userId, _ := middleware.GetUserID(c)

	data, insErr := ph.service.GetRevisionFromProject(c.Request.Context(), userId, projectId)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return

	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting revision data", data))
}

func (ph *ProjectHandler) PostCreateRevision(c *gin.Context) {

	param := c.Param("projectId")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	userId, _ := middleware.GetUserID(c)

	var req createRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, insErr := ph.service.CreateProjectRevision(c.Request.Context(), projectId, userId, req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully creating new revision", data))
}

func (ph *ProjectHandler) PatchUpdateRevisionStatus(c *gin.Context) {

	param := c.Param("projectId")
	projectId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}
	paramId := c.Param("id")
	revisionId, err := uuid.Parse(paramId)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid project id"))
		return
	}

	userId, _ := middleware.GetUserID(c)

	var req updateRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMap, isValid := pkg.ParseValidationErrors(err)

		if isValid {
			c.JSON(400, shared.NewErrorResponse(400, errMap))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, "invalid request body"))
		return
	}

	data, uptErr := ph.service.UpdateRevisionStatus(c.Request.Context(), userId, projectId, revisionId, req.Status)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updating revision status", data))
}

func (ph *ProjectHandler) RegisterRoutes(r gin.IRouter) {

	projectApi := r.Group("/project")
	projectApi.Use(middleware.AuthMiddleware())
	projectApi.GET("/details/:id", ph.HandleGetDetail)

	projectAdmin := projectApi.Group("/")
	projectAdmin.Use(middleware.RoleMiddleware(middleware.RoleAdmin))

	projectAdmin.GET("/get-all", ph.HandleGetAllProjects)
	projectAdmin.POST("/add", ph.PostCreateProject)
	projectAdmin.DELETE("/delete/:id", ph.DeleteProjectHandle)
	projectAdmin.PATCH("/update/:id", ph.PatchProjectData)

	// member client endpoint
	projectApi.GET("/:projectId/members", ph.HandleGetAllMember)
	projectApi.GET("/:projectId/progress", ph.GetProgressFromProject)

	// admin members endpoint
	projectAdmin.POST("/:projectId/members/add", ph.PostNewMember)
	projectAdmin.PATCH("/:projectId/members/update", ph.PatchMemberData)
	projectAdmin.DELETE("/:projectId/members/delete/:memberId", ph.DeleteProjectMember)

	// progress admin endpoint
	projectAdmin.POST("/:projectId/progress/add", ph.PostCreateProgress)
	projectAdmin.DELETE("/:projectId/progress/delete/:id", ph.DeleteProjectProgress)
	projectAdmin.PATCH("/:projectId/progress/update/:id", ph.PatchUpdateProgress)

	// ================== revision ====================
	projectApi.GET("/:projectId/revision", ph.GetRevisionFromProject)
	projectApi.POST("/:projectId/revision/add", ph.PostCreateRevision)
	projectAdmin.PATCH("/:projectId/revision/update/:id", ph.PatchUpdateRevisionStatus)
	// ================================================

}
