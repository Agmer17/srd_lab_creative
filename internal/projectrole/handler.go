package projectrole

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectRoleHandler struct {
	service *ProjectRoleService
}

func NewProjectRoleHandler(svc *ProjectRoleService) *ProjectRoleHandler {
	return &ProjectRoleHandler{
		service: svc,
	}
}

func (prh *ProjectRoleHandler) HandleGetAllRole(c *gin.Context) {

	data, err := prh.service.GetAllProjectRoles(c.Request.Context())
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting role data", data))
}

func (prh *ProjectRoleHandler) HandleSearchRole(c *gin.Context) {

	query, _ := c.GetQuery("q")
	data, err := prh.service.SearchRole(c.Request.Context(), query)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting role data", data))
}

func (prh *ProjectRoleHandler) PostCreateNewRole(c *gin.Context) {

	var req createRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid valid body for create role request"))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, insErr := prh.service.CreateRole(c.Request.Context(), req.Name)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully creating new role!", data))
}

func (prh *ProjectRoleHandler) PatchRole(c *gin.Context) {

	param := c.Param("id")
	roleId, err := uuid.Parse(param)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id parameter! please only using uuid for the id!"))
		return
	}

	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid valid body for update role request"))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, uptErr := prh.service.UpdateRole(c.Request.Context(), req.Name, roleId)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updated role data", data))
}

func (prh *ProjectRoleHandler) DeleteRoleHandler(c *gin.Context) {

	roleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id parameter! please only using uuid for the id!"))
		return
	}

	delErr := prh.service.DeleteRoles(c.Request.Context(), roleId)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully deleted role data", gin.H{
		"deleted_at": time.Now(),
	}))
}

func (prh *ProjectRoleHandler) RegisterRoutes(r gin.IRouter) {

	roleApi := r.Group("/project-role")

	roleApi.Use(middleware.AuthMiddleware())
	roleApi.Use(middleware.RoleMiddleware(middleware.RoleAdmin))

	roleApi.GET("/get-all", prh.HandleGetAllRole)
	roleApi.GET("/search", prh.HandleSearchRole)
	roleApi.POST("/add", prh.PostCreateNewRole)
	roleApi.PATCH("/update/:id", prh.PatchRole)
	roleApi.DELETE("/delete/:id")

}
